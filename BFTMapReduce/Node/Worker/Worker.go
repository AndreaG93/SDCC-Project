package Worker

import (
	"SDCC-Project/BFTMapReduce"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/BFTMapReduce/Task"
	"SDCC-Project/cloud/zookeeper"
	"encoding/gob"
	"fmt"
)

const ()

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
	(*output).mapReduceRPCAddress = fmt.Sprintf("%s:%d", internetAddress, BFTMapReduce.MapReduceRPCBasePort+id)
	(*output).mapReduceGetRPCAddress = fmt.Sprintf("%s:%d", internetAddress, BFTMapReduce.MapReduceGetRPCBasePort+id)
	(*output).zookeeperClient = zookeeper.New([]string{"localhost:2181"})

	return output
}

func (obj *Worker) StartWork() {

	gob.Register(WordCount.MapInput{})
	gob.Register(WordCount.ReduceInput{})

	go BFTMapReduce.StartAcceptingRPCRequest(&Task.MapReduce{}, (*obj).mapReduceRPCAddress)
	go BFTMapReduce.StartAcceptingRPCRequest(&Task.MapReduceGet{}, (*obj).mapReduceGetRPCAddress)

	(*obj).zookeeperClient.RegisterNodeMembership((*obj).id, (*obj).internetAddress)
	(*obj).zookeeperClient.KeepConnectionAlive()
}
