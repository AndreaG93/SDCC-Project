docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

docker run --name my-zookeeper --network host --restart always -d zookeeper

# docker exec -it my-zookeeper /bin/bash