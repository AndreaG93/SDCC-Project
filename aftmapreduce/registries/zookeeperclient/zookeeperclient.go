package zookeeperclient

import (
	"SDCC-Project/cloud/zookeeper"
	"sync"
)

var zookeeperClient *zookeeper.Client
var once sync.Once

func GetInstance() *zookeeper.Client {
	once.Do(func() {
		zookeeperClient = zookeeper.New([]string{"localhost:2181"})
	})
	return zookeeperClient
}
