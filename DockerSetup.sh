
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi -f $(docker images -q)


docker build --file DockerFile.ZooKeeper --tag zookeeper_image .
docker create --network host --name zookeeper_server zookeeper_image

docker container start zookeeper_server

# docker exec -it zookeeper_server /bin/bash
