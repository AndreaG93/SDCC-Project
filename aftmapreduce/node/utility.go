package node

import (
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/nodelogger"
)

var zookeeperClient *zookeeper.Client
var logger *nodelogger.Logger

func Initialize(nodeID int, nodeType string) {

	zookeeperClient = zookeeper.New([]string{"localhost:2181"})
	logger = nodelogger.New(nodeID, nodeType)
}

func GetZookeeperClient() *zookeeper.Client {
	return zookeeperClient
}

func GetLogger() *nodelogger.Logger {
	return logger
}
