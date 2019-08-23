go get -u github.com/Sirupsen/logrus
go get -u github.com/aws/aws-sdk-go/service/s3/...
go get -u github.com/aws/aws-sdk-go/aws/...
go get -u github.com/samuel/go-zookeeper/zk

docker run --name my-zookeeper --network host --restart always -d zookeeper
# docker exec -it my-zookeeper /bin/bash
# docker stop $(docker ps -a -q)
# docker rm $(docker ps -a -q)