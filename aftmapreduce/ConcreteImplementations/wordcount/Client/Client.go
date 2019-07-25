package Client

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/ConcreteImplementations/wordcount"
	"SDCC-Project/utility"
	"encoding/gob"
	"fmt"
	"net/rpc"
)

func SendRequest(digestFile string, internetAddress string) {

	input := new(aftmapreduce.MapReduceRequestInput)
	output := new(aftmapreduce.MapReduceRequestOutput)

	inputData := new(wordcount.Input)
	(*inputData).FileDigest = digestFile
	(*inputData).MapCardinality = 5

	input.InputData = inputData

	gob.Register(wordcount.Input{})

	currentLeaderPublicInternetAddress := fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRequestRPCBasePort+1)

	client, err := rpc.Dial("tcp", currentLeaderPublicInternetAddress)
	utility.CheckError(err)

	err = client.Call("MapReduceRequest.Execute", &input, &output)
	utility.CheckError(err)
}
