package Primary

import (
	"SDCC-Project/MapReduce"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/MapReduce/Task"
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
	gob.Register(WordCount.RawInput{})
	MapReduce.StartAcceptingRPCRequest(&Task.MapReduceRequest{}, (*obj).mapReduceRequestRPCAddress)
}
