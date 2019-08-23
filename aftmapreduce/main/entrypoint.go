package main

import (
	"SDCC-Project/aftmapreduce/node/primary"
	"SDCC-Project/aftmapreduce/node/worker"
	"SDCC-Project/aftmapreduce/utility"
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

	if nodeClass == "primary" {

		primary.Initialize(nodeID, nodePublicIP, zookeeperServers)
		primary.StartWork()

	} else if nodeClass == "worker" {

		worker.Initialize(nodeID, nodeGroupId, nodePublicIP, zookeeperServers)
		worker.StartWork()
	}
}
