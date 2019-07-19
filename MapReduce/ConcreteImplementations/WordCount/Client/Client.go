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

	inputData := new(WordCount.RawInput)
	(*inputData).FileDigest = digestFile
	(*inputData).MapCardinality = 5

	input.InputData = inputData

	gob.Register(WordCount.RawInput{})

	client, err := rpc.Dial("tcp", primaryNodeAddress)
	utility.CheckError(err)

	err = client.Call("MapReduceRequest.Execute", &input, &output)
	utility.CheckError(err)
}
