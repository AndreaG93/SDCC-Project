chmod 077 graziani-01.pem
index=0

ssh -i "graziani-01.pem" ubuntu@ec2-35-175-240-107.compute-1.amazonaws.com

# Setup EC2 instance servers with Apache-Zookeeper
ZkServersIPv4=("3.87.219.134" "3.94.62.19" "54.243.4.159")
for i in "${ZkServersIPv4[@]}"
do
    ((index++))
    ssh -i "graziani-01.pem" ubuntu@$i "
    
sudo apt update -y && sudo apt upgrade -y && sudo apt install -y golang



    "
done


# Start Apache-Zookeeper Servers
for i in "${ZkServersIPv4[@]}"
do
    ssh -i "graziani-01.pem" ubuntu@$i "sudo /usr/local/zookeeper/bin/zkServer.sh start"
done

