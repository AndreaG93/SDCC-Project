package main

import (
	"SDCC-Project/BFTMapReduce/Node/Worker"
	"SDCC-Project/utility"
	"os"
	"strconv"
)

func main() {

	id, err := strconv.Atoi(os.Getenv("NODE_ID"))
	utility.CheckError(err)

	Worker.New(id, string(30000+id), string(40000+id))
}
