# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
# Amazon EC2 instances configuration
#
# sudo apt install awscli jq (on Ubuntu)
# chmod u=xrw,g=,o= ./graziani.pem
#
# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# System's internet addesses...
PRIMARY_INSTENCES_PUBLIC_IP=()
PRIMARY_INSTENCES_PRIVATE_IP=()

ZOOKEEPER_INSTENCES_PRIVATE_IP=()
ZOOKEEPER_INSTENCES_PUBLIC_IP=()

WORKER_GROUP_0_INSTENCES_PRIVATE_IP=()
WORKER_GROUP_0_INSTENCES_PUBLIC_IP=()

WORKER_GROUP_1_INSTENCES_PRIVATE_IP=()
WORKER_GROUP_1_INSTENCES_PUBLIC_IP=()

# ======================================================================================================================
# Zookeeper Servers
# ======================================================================================================================

# 1 - Information retrival...
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")

    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].PrivateIpAddress')
    OUTPUT2=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ZOOKEEPER_INSTENCES_PRIVATE_IP+=("$OUTPUT1")
    ZOOKEEPER_INSTENCES_PUBLIC_IP+=("$OUTPUT2")
done

# 3 - Start
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    ssh -i "graziani.pem" ubuntu@$IP "sudo /usr/local/zookeeper/bin/zkServer.sh start"
done

# 4 - Status
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    ssh -i "graziani.pem" ubuntu@$IP "sudo /usr/local/zookeeper/bin/zkServer.sh status"
done

# ======================================================================================================================
# Primary Servers
# ======================================================================================================================

# Information retrival...
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")
    OUTPUT=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    PRIMARY_INSTENCES_PUBLIC_IP+=("$OUTPUT")
done

# Configuration
index=0
for i in "${PRIMARY_INSTENCES_PUBLIC_IP[@]}"
do
    ((index++))
    ssh -i "graziani.pem" ubuntu@$i "

sudo apt update && sudo apt install -y golang jq

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk
go get -u github.com/Sirupsen/logrus

rm -R ./go/src/SDCC-Project
git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

jq -n '{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[0]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[1]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[2]}\"], NodeID: $index, NodeGroupID: 0, NodeClass: \"Primary\"}' > ./conf.json"
done

# 3 - Start
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani.pem" ubuntu@$IP "
    go run ./go/src/SDCC-Project/aftmapreduce/main/entrypoint.go &"
done

# ======================================================================================================================
# Worker Servers GROUP 0
# ======================================================================================================================

# Information retrival...
for i in 0 1 2
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$id","Name=tag:ID-Group,Values=0")
    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    WORKER_GROUP_0_INSTENCES_PUBLIC_IP+=("$OUTPUT1")
done

# Configuration
index=0
for i in "${WORKER_GROUP_0_INSTENCES_PUBLIC_IP[@]}"
do
    ssh -i "graziani.pem" ubuntu@$i "

sudo apt update && sudo apt install -y golang jq sysstat

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk
go get -u github.com/Sirupsen/logrus

rm -R ./go/src/SDCC-Project
git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

jq -n '{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[0]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[1]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[2]}\"], NodeID: $index, NodeGroupID: 0, NodeClass: \"Worker\"}' > ./conf.json"
  ((index++))
done

# 3 - Start
for i in 0 1 2
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$id","Name=tag:ID-Group,Values=0")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani.pem" ubuntu@$IP "
    go run ./go/src/SDCC-Project/aftmapreduce/main/entrypoint.go &"
done

# ======================================================================================================================
# Worker Servers GROUP 1
# ======================================================================================================================

# Information retrival...
for i in 3 4 5
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$id","Name=tag:ID-Group,Values=1")
    OUTPUT1=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    WORKER_GROUP_0_INSTENCES_PUBLIC_IP+=("$OUTPUT1")
done

# Configuration
index=3
for i in "${WORKER_GROUP_0_INSTENCES_PUBLIC_IP[@]}"
do
    ssh -i "graziani.pem" ubuntu@$i "

sudo apt update && sudo apt install -y golang jq

go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk
go get -u github.com/Sirupsen/logrus

rm -R ./go/src/SDCC-Project
git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project

jq -n '{ZookeeperServersPrivateIPs: [\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[0]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[1]}\",\"${ZOOKEEPER_INSTENCES_PRIVATE_IP[2]}\"], NodeID: $index, NodeGroupID: 0, NodeClass: \"Worker\"}' > ./conf.json"
    ((index++))
done

# 3 - Start
for i in 3 4 5
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$id","Name=tag:ID-Group,Values=1")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani.pem" ubuntu@$IP "
    go run ./go/src/SDCC-Project/aftmapreduce/main/entrypoint.go &"
done









































