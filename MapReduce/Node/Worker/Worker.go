package Worker

import (
	"SDCC-Project/MapReduce"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/MapReduce/Task"
	"encoding/gob"
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

	gob.Register(WordCount.MapInput{})
	gob.Register(WordCount.ReduceInput{})

	go MapReduce.StartAcceptingRPCRequest(&Task.MapReduce{}, (*obj).mapReduceRPCAddress)
	MapReduce.StartAcceptingRPCRequest(&Task.MapReduceGet{}, (*obj).mapReduceGetRPCAddress)
}
