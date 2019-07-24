package Client

import (
	"SDCC-Project/BFTMapReduce"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/BFTMapReduce/Task"
	"SDCC-Project/utility"
	"encoding/gob"
	"fmt"
	"net/rpc"
)

func SendRequest(digestFile string, internetAddress string) {

	input := new(Task.MapReduceRequestInput)
	output := new(Task.MapReduceRequestOutput)

	inputData := new(WordCount.File)
	(*inputData).FileDigest = digestFile
	(*inputData).MapCardinality = 5

	input.InputData = inputData

	gob.Register(WordCount.File{})

	currentLeaderPublicInternetAddress := fmt.Sprintf("%s:%d", internetAddress, BFTMapReduce.MapReduceRequestRPCBasePort+1)

	client, err := rpc.Dial("tcp", currentLeaderPublicInternetAddress)
	utility.CheckError(err)

	err = client.Call("MapReduceRequest.Execute", &input, &output)
	utility.CheckError(err)
}
