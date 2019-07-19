package Worker

import (
	"SDCC-Project/MapReduce"
	"SDCC-Project/MapReduce/Task"
)

type Worker struct {
	id                     int
	mapReduceRPCAddress    string
	mapReduceGetRPCAddress string
}

func New(id int, mapReduceRPCAddress string, mapReduceGetRPCAddress string) *Worker {

	output := new(Worker)
	(*output).id = id
	(*output).mapReduceRPCAddress = mapReduceRPCAddress
	(*output).mapReduceGetRPCAddress = mapReduceGetRPCAddress

	return output
}

func (obj *Worker) StartWork() {

	go MapReduce.StartAcceptingRPCRequest(Task.MapReduce{}, (*obj).mapReduceRPCAddress)
	MapReduce.StartAcceptingRPCRequest(Task.MapReduceGet{}, (*obj).mapReduceGetRPCAddress)
}
