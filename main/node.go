package main

import (
	"SDCC-Project/aftmapreduce/nodes/primary"
	"SDCC-Project/aftmapreduce/nodes/worker"
	"SDCC-Project/utility"
	"os"
	"strconv"
)

func main() {

	nodeID, err := strconv.Atoi(os.Args[1])
	utility.CheckError(err)

	nodeClass := os.Args[2]
	nodePublicIP := os.Args[3]

	if nodeClass == "primary" {

		primary.New(nodeID, nodePublicIP).StartWork()

	} else if nodeClass == "worker" {

		worker.New(nodeID, nodePublicIP).StartWork()
	}
}
