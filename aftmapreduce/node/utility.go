package node

import (
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/nodelogger"
)

var zookeeperClient *zookeeper.Client
var logger *nodelogger.Logger

func Initialize(nodeID int, nodeType string, zookeeperAddresses []string) {

	zookeeperClient = zookeeper.New(zookeeperAddresses)
	logger = nodelogger.New(nodeID, nodeType)
}

func GetZookeeperClient() *zookeeper.Client {
	return zookeeperClient
}

func GetLogger() *nodelogger.Logger {
	return logger
}
