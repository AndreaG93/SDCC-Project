package main

import (
	"SDCC-Project/aftmapreduce/process/client"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("USAGE: %s [INPUT FILE PATH] [Zookeeper-Server-1 IPv4] [Zookeeper-Server-2 IPv4]...", os.Args[0])
	}

	inputFilePath := os.Args[1]
	zookeeperServerInternetAddresses := os.Args[2:]

	client.StartWork(inputFilePath, zookeeperServerInternetAddresses)
}
