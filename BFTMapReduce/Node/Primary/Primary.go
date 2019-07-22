package Primary

import (
	"SDCC-Project/BFTMapReduce"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/BFTMapReduce/Task"
	"encoding/gob"
)

type Primary struct {
	id                         int
	mapReduceRequestRPCAddress string
}

func New(id int, mapReduceRequestRPCAddress string) *Primary {

	output := new(Primary)
	(*output).id = id
	(*output).mapReduceRequestRPCAddress = mapReduceRequestRPCAddress

	return output
}

func (obj *Primary) StartWork() {
	gob.Register(WordCount.File{})
	gob.Register(WordCount.MapInput{})
	gob.Register(WordCount.ReduceInput{})

	BFTMapReduce.StartAcceptingRPCRequest(&Task.MapReduceRequest{}, (*obj).mapReduceRequestRPCAddress)
}
