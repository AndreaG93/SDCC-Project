sudo apt install python-pip
sudo pip install awscli

# ==================================================================================================================== #
# Zookeeper EC2 instances data retrival
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

# ==================================================================================================================== #
# Zookeeper Servers configuration
# ==================================================================================================================== #

index=0
for i in "${ZOOKEEPER_SERVER_PUBLIC_IP[@]}"
do
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

# ==================================================================================================================== #
# Primary Servers configuration
# ==================================================================================================================== #












jq -n "{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_SERVER_PRIVATE_IP[0]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[1]}\",\"${ZOOKEEPER_SERVER_PRIVATE_IP[2]}\"], NodeID: 1, NodeClass: \"Worker\"}" > conf.json


