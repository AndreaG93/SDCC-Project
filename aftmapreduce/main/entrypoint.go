package main

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/primary"
	"SDCC-Project/aftmapreduce/process/worker"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"os/exec"
)

func main() {

	configuration := new(utility.NodeConfiguration)
	configuration.Load("./conf.json")

	command := exec.Command("curl", "http://169.254.169.254/latest/meta-data/public-ipv4")
	commandOutput, err := command.Output()
	if err != nil {
		panic(err)
	}

	nodeID := configuration.NodeID
	nodeClass := configuration.NodeClass
	nodeGroupId := configuration.NodeGroupID
	zookeeperServers := configuration.ZookeeperServersPrivateIPs
	nodePublicIP := string(commandOutput)

	utility.CheckError(process.Initialize(nodeID, nodeGroupId, nodeClass, nodePublicIP, zookeeperServers))

	if nodeClass == "Primary" {
		fmt.Printf("Start as Primary -- Node %d", nodeID)
		primary.StartWork()

	} else if nodeClass == "Worker" {
		fmt.Printf("Start as Worker -- Node %d, Group %d", nodeID, nodeGroupId)
		worker.StartWork()
	}
}
