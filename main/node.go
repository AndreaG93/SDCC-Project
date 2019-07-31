package main

import (
	"SDCC-Project/aftmapreduce/node/primary"
	"SDCC-Project/aftmapreduce/node/worker"
	"SDCC-Project/utility"
	"os/exec"
)

func main() {

	configuration := new(utility.NodeConfiguration)
	configuration.Load("conf.json")

	command := exec.Command("curl", "http://169.254.169.254/latest/meta-data/public-ipv4")
	commandOutput, err := command.Output()
	if err != nil {
		panic(err)
	}

	nodeID := configuration.NodeID
	nodeClass := configuration.NodeClass
	nodePublicIP := string(commandOutput)

	if nodeClass == "primary" {

		primary.New(nodeID, nodePublicIP).StartWork()

	} else if nodeClass == "worker" {

		worker.New(nodeID, nodePublicIP).StartWork()
	}
}
