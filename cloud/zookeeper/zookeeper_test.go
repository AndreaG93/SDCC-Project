package zookeeper

import (
	"strings"
	"testing"
)

const (
	testData      = "test"
	testZNodePath = "/test"
)

func Test_zookeeperBasicOperations(t *testing.T) {

	zooKeeperClient := New([]string{"localhost:2181"})

	if !(*zooKeeperClient).CheckZNodeExistence(testZNodePath) {
		(*zooKeeperClient).CreateZNode(testZNodePath)
	}

	(*zooKeeperClient).SetZNodeData(testZNodePath, []byte(testData))

	dataReceived, _ := (*zooKeeperClient).GetZNodeData(testZNodePath)

	if strings.Compare(testData, string(dataReceived)) != 0 {
		panic("Error")
	}

	(*zooKeeperClient).RemoveZNode(testZNodePath)
}
