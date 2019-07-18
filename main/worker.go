package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/node/worker"
	"SDCC-Project-WorkerNode/utility"
	"os"
	"strconv"
)

func main() {

	myId, err := strconv.Atoi(os.Getenv("NODE_ID"))
	utility.CheckError(err)

	node := worker.New(uint(myId), []string{"localhost"})
	node.StartWork()
}
