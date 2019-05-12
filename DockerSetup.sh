
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
# docker rmi -f $(docker images -q)


docker build --file DockerFile.ZooKeeper --tag zookeeper_image .
docker build --file DockerFile.PrimaryNode --tag primary_node .

docker create --network host --name zookeeper_server zookeeper_image

docker container start zookeeper_server

# Primary nodes...

docker network create -d bridge mybridge

docker create --network=host --name primary1 --env NODE_ID=1 primary_node
docker create --network=host --name primary2 --env NODE_ID=2 primary_node
docker create --network=host --name primary3 --env NODE_ID=3 primary_node
docker create --network=host --name primary4 --env NODE_ID=4 primary_node
docker create --network=host --name primary5 --env NODE_ID=5 primary_node

docker container start primary1
docker container start primary2
docker container start primary3
docker container start primary4
docker container start primary5

# docker exec -it zookeeper_server /bin/bash
# docker exec -it primary1 /go/primarynode
# zkCli.sh
