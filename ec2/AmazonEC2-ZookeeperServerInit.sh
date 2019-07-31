chmod 077 graziani-01.pem
index=0

# Setup EC2 instance servers with Apache-Zookeeper
ZkServersIPv4=("3.87.219.134" "3.94.62.19" "54.243.4.159")
for i in "${ZkServersIPv4[@]}"
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
server.1=172.31.83.14:2888:3888
server.2=172.31.92.183:2888:3888
server.3=172.31.85.212:2888:3888' | sudo tee /usr/local/zookeeper/conf/zoo.cfg

echo $index | sudo tee /var/lib/zookeeper/myid    
    "
done


# Start Apache-Zookeeper Servers
for i in "${ZkServersIPv4[@]}"
do
    ssh -i "graziani-01.pem" ubuntu@$i "sudo /usr/local/zookeeper/bin/zkServer.sh start"
done

# Check status of Apache-Zookeeper Servers
for i in "${ZkServersIPv4[@]}"
do
    ssh -i "graziani-01.pem" ubuntu@$i "sudo /usr/local/zookeeper/bin/zkServer.sh status"
done
