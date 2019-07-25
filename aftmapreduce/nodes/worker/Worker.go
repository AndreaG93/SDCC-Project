package worker

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/ConcreteImplementations/wordcount"
	"SDCC-Project/cloud/zookeeper"
	"encoding/gob"
	"fmt"
)

type Worker struct {
	id                     int
	internetAddress        string
	mapReduceRPCAddress    string
	mapReduceGetRPCAddress string
	zookeeperClient        *zookeeper.Client
}

func New(id int, internetAddress string) *Worker {

	output := new(Worker)

	(*output).id = id
	(*output).internetAddress = internetAddress
	(*output).mapReduceRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRPCBasePort+id)
	(*output).mapReduceGetRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceGetRPCBasePort+id)
	(*output).zookeeperClient = zookeeper.New([]string{"localhost:2181"})

	return output
}

func (obj *Worker) StartWork() {

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	go aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.Replica{}, (*obj).mapReduceRPCAddress)
	go aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.DataRetriever{}, (*obj).mapReduceGetRPCAddress)

	(*obj).zookeeperClient.RegisterNodeMembership((*obj).id, (*obj).internetAddress)
	(*obj).zookeeperClient.KeepConnectionAlive()
}
