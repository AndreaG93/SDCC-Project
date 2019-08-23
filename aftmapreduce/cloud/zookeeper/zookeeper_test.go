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

	zooKeeperClient := New([]string{"3.87.219.134:2181", "3.94.62.19:2181", "54.243.4.159:2181"})
	if !(*zooKeeperClient).CheckZNodeExistence(testZNodePath) {
		(*zooKeeperClient).CreateZNode(testZNodePath, nil, int32(0))
	}

	(*zooKeeperClient).SetZNodeData(testZNodePath, []byte(testData))

	dataReceived, _ := (*zooKeeperClient).GetZNodeData(testZNodePath)

	if strings.Compare(testData, string(dataReceived)) != 0 {
		panic("Error")
	}

	(*zooKeeperClient).RemoveZNode(testZNodePath)
}
