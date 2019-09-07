sudo apt install python-pip
sudo pip install awscli
chmod 077 graziani-01.pem

# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
# Zookeeper nodes configuration on Amazon EC2
# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# Information retrival...
# ==================================================================================================================== #
ZOOKEEPER_SERVER_PRIVATE_IP=()
ZOOKEEPER_SERVER_PUBLIC_IP=()

for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")

    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].PrivateIpAddress')
    OUTPUT2=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ZOOKEEPER_SERVER_PRIVATE_IP+=("$OUTPUT1")
    ZOOKEEPER_SERVER_PUBLIC_IP+=("$OUTPUT2")
done

# Configuration
# ==================================================================================================================== #

index=0
for i in "${ZOOKEEPER_SERVER_PUBLIC_IP[@]}"
do
    ((index++))
    ssh -i "graziani-01.pem" ubuntu@$i "

sudo apt update -y && sudo apt upgrade -y && sudo apt install -y default-jre

sudo wget https://www-us.apache.org/dist/zookeeper/zookeeper-3.5.5/apache-zookeeper-3.5.5-bin.tar.gz

sudo tar -xzf  apache-zookeeper-3.5.5-bin.tar.gz
sudo mv apache-zookeeper-3.5.5-bin /usr/local/zookeeper

echo 'tickTime=2000
initLimit=10
syncLimit=5
dataDir=/var/lib/zookeeper
clientPort=2181
server.1=${ZOOKEEPER_SERVER_PRIVATE_IP[0]}:2888:3888
server.2=${ZOOKEEPER_SERVER_PRIVATE_IP[1]}:2888:3888
server.3=${ZOOKEEPER_SERVER_PRIVATE_IP[2]}:2888:3888' | sudo tee /usr/local/zookeeper/conf/zoo.cfg

echo $index | sudo tee /var/lib/zookeeper/myid
    "
done

# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
# Primary nodes configuration on Amazon EC2
# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# Information retrival...
# ==================================================================================================================== #

PRIMARY_SERVERS_PUBLIC_IP=()

for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=PrimaryServer" "Name=tag:ID,Values=$i")

    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    PRIMARY_SERVERS_PUBLIC_IP+=("$OUTPUT1")
done

# Configuration
# ==================================================================================================================== #

index=0
for i in "${PRIMARY_SERVERS_PUBLIC_IP[@]}"
do
    ((index++))
    ssh -i "graziani-01.pem" ubuntu@$i "

sudo apt update -y && sudo apt upgrade -y

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk
go get -u github.com/Sirupsen/logrus

git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

jq -n "{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_SERVER_PRIVATE_IP[0]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[1]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[2]}\"], NodeID: $index, NodeGroupID: 0, NodeClass: \"Primary\"}" > ./go/src/SDCC-Project/main/conf.json"
done

# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
# Worker nodes configuration on Amazon EC2
# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# Information retrival...
# ==================================================================================================================== #

WORKER_SERVERS_PUBLIC_IP()

for id in 0 1 2 3 4 5 6 7 8
for groupid in 0 1 2
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$id","Name=tag:ID-Group,Values=$groupid")

    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    WORKER_SERVERS_PUBLIC_IP+=("$OUTPUT1")
done
done

# Configuration
# ==================================================================================================================== #

index=0
for i in "${WORKER_SERVERS_PUBLIC_IP[@]}"
do
    ((index++))
    ssh -i "graziani-01.pem" ubuntu@$i "

sudo apt update -y && sudo apt upgrade -y

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/z

git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

jq -n "{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_SERVER_PRIVATE_IP[0]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[1]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[2]}\"], NodeID: $index, NodeClass: \"Worker\"}" > ./go/src/SDCC-Project/main/conf.json

    "
done
