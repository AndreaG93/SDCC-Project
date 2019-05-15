package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
)

func main() {

	var leaderAddress string
	var jobWatcher <-chan zk.Event

	leaderAddress, jobWatcher = zookeeper.GetActualClusterLeaderAddress()
	fmt.Println(leaderAddress)

	<-jobWatcher

	fmt.Println(leaderAddress)
}
