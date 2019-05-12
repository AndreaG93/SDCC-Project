package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/primarynode"
	"SDCC-Project-WorkerNode/utility"
	"os"
	"strconv"
)

func main() {

	myId, err := strconv.Atoi(os.Getenv("NODE_ID"))
	utility.CheckError(err)

	node := primarynode.New(uint(myId))
	node.StartWork()
}
