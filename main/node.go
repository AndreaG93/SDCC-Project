package main

import (
	"SDCC-Project/BFTMapReduce/Node/Primary"
	"SDCC-Project/BFTMapReduce/Node/Worker"
	"SDCC-Project/utility"
	"os"
	"strconv"
)

func main() {

	nodeID, err := strconv.Atoi(os.Args[1])
	utility.CheckError(err)

	nodeClass := os.Args[2]
	nodePublicIP := os.Args[3]

	if nodeClass == "Primary" {

		Primary.New(nodeID, nodePublicIP).StartWork()

	} else if nodeClass == "Worker" {

		Worker.New(nodeID, nodePublicIP).StartWork()
	}
}
