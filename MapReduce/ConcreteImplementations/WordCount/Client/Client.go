package Client

import (
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/MapReduce/Task"
	"SDCC-Project/utility"
	"encoding/gob"
	"net/rpc"
)

func SendRequest(digestFile string, primaryNodeAddress string) {

	input := new(Task.MapReduceRequestInput)
	output := new(Task.MapReduceRequestOutput)

	inputData := new(WordCount.File)
	(*inputData).FileDigest = digestFile
	(*inputData).MapCardinality = 5

	input.InputData = inputData

	gob.Register(WordCount.File{})

	client, err := rpc.Dial("tcp", primaryNodeAddress)
	utility.CheckError(err)

	err = client.Call("MapReduceRequest.Execute", &input, &output)
	utility.CheckError(err)
}
