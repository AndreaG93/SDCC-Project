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

	input := new(aftmapreduce.EntryPointInput)
	output := new(aftmapreduce.EntryPointOutput)

	inputData := new(wordcount.Input)
	(*inputData).FileDigest = "test"
	(*inputData).MapCardinality = 5

	input.Data = inputData

	gob.Register(wordcount.Input{})

	currentLeaderPublicInternetAddress := fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRequestRPCBasePort+1)

	client, err := rpc.Dial("tcp", currentLeaderPublicInternetAddress)
	utility.CheckError(err)

	err = client.Call("EntryPoint.Execute", &input, &output)
	utility.CheckError(err)
}
