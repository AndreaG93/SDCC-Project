PP_INSTENCES_PUBLIC_IP=()
PP_INSTENCES_PRIVATE_IP=()

WP_INSTENCES_PUBLIC_IP=()
WP_INSTENCES_PRIVATE_IP=()

# ==================================================================================================================== #
# Information retrival...
# ==================================================================================================================== #
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Name,Values=PP-$i")
  
    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    OUTPUT2=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].PrivateIpAddress')
   
    PP_INSTENCES_PUBLIC_IP+=("$OUTPUT1")
    PP_INSTENCES_PRIVATE_IP+=("$OUTPUT2")
done

for i in 0 1 2 3 4 5
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Name,Values=WP-$i")
  
    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    OUTPUT2=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].PrivateIpAddress')
   
    WP_INSTENCES_PUBLIC_IP+=("$OUTPUT1")
    WP_INSTENCES_PRIVATE_IP+=("$OUTPUT2")
done

# To print IPs for test use: echo ${PP_INSTENCES_PUBLIC_IP[*]}
# To print IPs for test use: echo ${WP_INSTENCES_PUBLIC_IP[*]}

# ==================================================================================================================== #
# Build application on all AWS instances...
# ==================================================================================================================== #

INSTENCES=( "${PP_INSTENCES_PUBLIC_IP[@]}" "${WP_INSTENCES_PUBLIC_IP[@]}" )

for i in ${INSTENCES[@]}
do

echo "Connecting To ... $i"

konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@$i "

sudo apt update && sudo apt install -y golang sysstat

echo DownloadGoData

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk
go get -u github.com/Sirupsen/logrus

rm -R ./go/src/SDCC-Project
git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

go build -o ./sddc_wc_aws_build ./go/src/SDCC-Project/aftmapreduce/main/wcserver.go

echo Complete
" &
done

# ==================================================================================================================== #
# Configuration Zookeeper cluster
# ==================================================================================================================== #

# -- tickTime

# the length of a single tick, which is the basic time unit used by ZooKeeper, as measured in milliseconds. It is used to regulate heartbeats, and timeouts. For example, # the minimum session timeout will be two ticks.

index=0
for i in ${PP_INSTENCES_PUBLIC_IP[@]}
do

echo "Connecting To ... $i"
((index++))
konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@$i "

sudo apt update && sudo apt install -y default-jre

sudo wget https://www-us.apache.org/dist/zookeeper/zookeeper-3.5.5/apache-zookeeper-3.5.5-bin.tar.gz

sudo tar -xzf  apache-zookeeper-3.5.5-bin.tar.gz
sudo rm -rf /usr/local/zookeeper
sudo mv apache-zookeeper-3.5.5-bin /usr/local/zookeeper

echo 'tickTime=50
initLimit=10
syncLimit=5
dataDir=/var/lib/zookeeper
clientPort=2181
server.1=${PP_INSTENCES_PRIVATE_IP[0]}:2888:3888
server.2=${PP_INSTENCES_PRIVATE_IP[1]}:2888:3888
server.3=${PP_INSTENCES_PRIVATE_IP[2]}:2888:3888' | sudo tee /usr/local/zookeeper/conf/zoo.cfg

sudo sh -c 'echo '$index' > /var/lib/zookeeper/myid'

echo Complete
" &
done

# ==================================================================================================================== #
# Generate all conf.json...
# ==================================================================================================================== #

index=0
for x in 0 1 2
do
    ((index++))
    MSG='{"ZookeeperServersPrivateIPs": ["'${PP_INSTENCES_PRIVATE_IP[0]}'","'${PP_INSTENCES_PRIVATE_IP[1]}'","'${PP_INSTENCES_PRIVATE_IP[2]}'"], "NodeID": '$index', "NodeGroupID": 0, "NodeClass": "Primary"}'

    echo $MSG | ssh -i "graziani.pem" ubuntu@${PP_INSTENCES_PUBLIC_IP[$x]} 'cat > ./conf.json'

done


for x in 0 1 2
do
    ((index++))
    MSG='{"ZookeeperServersPrivateIPs": ["'${PP_INSTENCES_PRIVATE_IP[0]}'","'${PP_INSTENCES_PRIVATE_IP[1]}'","'${PP_INSTENCES_PRIVATE_IP[2]}'"], "NodeID": '$x', "NodeGroupID": 0, "NodeClass": "Worker"}'

    echo $MSG | ssh -i "graziani.pem" ubuntu@${WP_INSTENCES_PUBLIC_IP[$x]} 'cat > ./conf.json'

done


for x in 3 4 5
do

    MSG='{"ZookeeperServersPrivateIPs": ["'${PP_INSTENCES_PRIVATE_IP[0]}'","'${PP_INSTENCES_PRIVATE_IP[1]}'","'${PP_INSTENCES_PRIVATE_IP[2]}'"], "NodeID": '$x', "NodeGroupID": 1, "NodeClass": "Worker"}'

    echo $MSG | ssh -i "graziani.pem" ubuntu@${WP_INSTENCES_PUBLIC_IP[$x]} 'cat > ./conf.json'

done


# ==================================================================================================================== #
# To start all processes...
# ==================================================================================================================== #


for x in ${PP_INSTENCES_PUBLIC_IP[@]}
do
    konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@$x 'sudo /usr/local/zookeeper/bin/zkServer.sh start' &
done

for x in ${PP_INSTENCES_PUBLIC_IP[@]}
do
    konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@$x 'sudo /usr/local/zookeeper/bin/zkServer.sh status' &
done

konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@${PP_INSTENCES_PUBLIC_IP[0]} './sddc_wc_aws_build' &

konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@${PP_INSTENCES_PUBLIC_IP[1]} './sddc_wc_aws_build' &
konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@${PP_INSTENCES_PUBLIC_IP[2]} './sddc_wc_aws_build' &

for x in ${WP_INSTENCES_PUBLIC_IP[@]}
do
    konsole --new-tab --noclose -e ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@$x './sddc_wc_aws_build' &
done

# ==================================================================================================================== #
# Local process...
# ==================================================================================================================== #

# ssh -o "StrictHostKeyChecking=no" -i "graziani.pem" ubuntu@3.229.132.166 '/usr/local/zookeeper/bin/zkCli.sh'

go build -o ./wcclient $HOME/go/src/SDCC-Project/aftmapreduce/process/client/main/wcclient.go
cp $HOME/go/src/SDCC-Project/test-input-data/input1.txt $HOME/input1.txt

./wcclient input1.txt 3.229.132.166 34.231.249.26 34.231.249.26



