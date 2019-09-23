package utility

import (
	"fmt"
	"testing"
)

func Test_Configuration(t *testing.T) {

	configuration := new(NodeConfiguration)
	if err := configuration.Load("conf.json"); err != nil {
		panic(err)
	}

	fmt.Println(configuration.ZookeeperServersPrivateIPs)
	fmt.Println(configuration.NodeClass)
	fmt.Println(configuration.NodeID)
	fmt.Println(configuration.NodeGroupID)
}
