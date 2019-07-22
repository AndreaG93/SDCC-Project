package zookeeper

import (
	"fmt"
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
		(*zooKeeperClient).CreateZNode(testZNodePath, 0)
	}

	(*zooKeeperClient).SetZNodeData(testZNodePath, []byte(testData))

	dataReceived, _ := (*zooKeeperClient).GetZNodeData(testZNodePath)

	if strings.Compare(testData, string(dataReceived)) != 0 {
		panic("Error")
	}

	(*zooKeeperClient).RemoveZNode(testZNodePath)
}

func Test_membershipWatcher(t *testing.T) {

	zooKeeperClient := New([]string{"localhost:2181"})

	for {
		data, watcher := zooKeeperClient.getMembershipZNodeData()
		fmt.Println(data)
		<-watcher
		fmt.Println("There are some changes...")
	}
}

func Test_ephemeralNodes(t *testing.T) {

	zooKeeperClient := New([]string{"localhost:2181"})
	(*zooKeeperClient).registerNodeMembership(1)

	_, channel := (*zooKeeperClient).GetZNodeData("/")
	<-channel
}
