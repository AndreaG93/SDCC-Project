package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"time"
)

func main() {

	system.SendHeartbeatTo("primary1", "primary2")

	time.Sleep(500 * time.Second)

	//network.StopHeartBeat()

}
