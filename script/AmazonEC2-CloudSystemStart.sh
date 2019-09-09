# Zookeeper Servers...
# ==================================================================================================================== #
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")
    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')
    ssh -i "graziani.pem" ubuntu@$IP "sudo /usr/local/zookeeper/bin/zkServer.sh start"
done

# Primary Servers...
# ==================================================================================================================== #
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=PrimaryServer" "Name=tag:ID,Values=$i")

    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani-01.pem" ubuntu@$IP "go run ./go/src/SDDC-Project/main/node.go"
done

# Worker Servers...
# ==================================================================================================================== #
for i in 1 2 3 4 5
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=Worker" "Name=tag:ID,Values=$i")

    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani-01.pem" ubuntu@$IP "go run ./go/src/SDDC-Project/main/node.go"
done



# ==================================================================================================================== #
# Check status of Apache-Zookeeper Servers
# ==================================================================================================================== #
for i in 1 2 3
do
    EC2_OUTPUT=$(aws ec2 describe-instances --region us-east-1 --filters "Name=tag:Role,Values=ZookeeperServer" "Name=tag:ID,Values=$i")

    IP=$(echo $EC2_OUTPUT | jq -r '.Reservations[].Instances[].NetworkInterfaces[].Association.PublicIp')

    ssh -i "graziani-01.pem" ubuntu@$IP "sudo /usr/local/zookeeper/bin/zkServer.sh status"
done
