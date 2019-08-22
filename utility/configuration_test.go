package utility

import (
	"fmt"
	"testing"
)

func Test_Configuration(t *testing.T) {

	configuration := new(NodeConfiguration)
	configuration.Load("example_conf.json")

	fmt.Println(configuration.ZookeeperServersPrivateIPs)
	fmt.Println(configuration.NodeClass)
	fmt.Println(configuration.NodeID)
	fmt.Println(configuration.NodeGroupID)
}
