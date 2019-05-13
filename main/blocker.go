package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	leaderOfflineChannel := make(chan os.Signal, 1)
	signal.Notify(leaderOfflineChannel, syscall.SIGUSR2)

	system.ReceiveHeartbeatFromSingleNode("primary2")

	<-leaderOfflineChannel

	fmt.Println("OK")

	/*
		i := 3

		for {

			i++
			time.Sleep(50 * time.Second)
			if i > 100 {
				i--
			}

		}
	*/

}
