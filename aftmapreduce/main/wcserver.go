package main

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/primary"
	"SDCC-Project/aftmapreduce/process/worker"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	var nodePublicIP string

	configuration := new(utility.NodeConfiguration)
	if err := configuration.Load("conf.json"); err != nil {
		panic(err)
	}

	if len(os.Args) == 1 {

		command := exec.Command("curl", "http://169.254.169.254/latest/meta-data/public-ipv4")
		if commandOutput, err := command.Output(); err != nil {
			panic(err)
		} else {
			nodePublicIP = string(commandOutput)
		}

	} else {
		nodePublicIP = "localhost"
	}

	nodeID := configuration.NodeID
	nodeClass := configuration.NodeClass
	nodeGroupId := configuration.NodeGroupID
	zookeeperServers := configuration.ZookeeperServersPrivateIPs

	utility.CheckError(process.Initialize(nodeID, nodeGroupId, nodeClass, nodePublicIP, zookeeperServers))

	if nodeClass == "Primary" {
		fmt.Printf("Start as Primary -- Node %d", nodeID)
		primary.StartWork()

	} else if nodeClass == "Worker" {
		fmt.Printf("Start as Worker -- Node %d, Group %d", nodeID, nodeGroupId)
		worker.StartWork()
	}
}
