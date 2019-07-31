package worker

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/implementations/wordcount"
	"SDCC-Project/aftmapreduce/node"
	"encoding/gob"
	"fmt"
)

type Worker struct {
	id                     int
	internetAddress        string
	mapReduceRPCAddress    string
	mapReduceGetRPCAddress string
}

func New(id int, internetAddress string, zookeeperAddresses []string) *Worker {

	output := new(Worker)

	node.Initialize(id, "Worker", zookeeperAddresses)

	(*output).id = id
	(*output).internetAddress = internetAddress
	(*output).mapReduceRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRPCBasePort+id)
	(*output).mapReduceGetRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceGetRPCBasePort+id)

	return output
}

func (obj *Worker) StartWork() {

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	go aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.Replica{}, (*obj).mapReduceRPCAddress)
	go aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.DataRetriever{}, (*obj).mapReduceGetRPCAddress)

	node.GetZookeeperClient().RegisterNodeMembership((*obj).id, (*obj).internetAddress)
	node.GetZookeeperClient().KeepConnectionAlive()
}
