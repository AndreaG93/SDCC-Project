package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"errors"
	"fmt"
	"net/rpc"
)

type Retrieve struct {
}

type RetrieveInput struct {
	DataDigest string
}

type RetrieveOutput struct {
	RawData []byte
}

func (x *Retrieve) Execute(input RetrieveInput, output *RetrieveOutput) error {

	process.GetLogger().PrintInfoTaskMessage(RetrieveTaskName, fmt.Sprintf("Received a 'RETRIEVE' request -- Data digest requested is %s", input.DataDigest))

	output.RawData = process.GetDataRegistry().Get(input.DataDigest)
	if output.RawData == nil {
		return errors.New("no data with given digest")
	}

	return nil
}

func retrieveFrom(NodeIPs []string, dataDigest string) []byte {

	process.GetLogger().PrintInfoTaskMessage(RetrieveTaskName, fmt.Sprintf("Target Nodes are %s", NodeIPs))

	var input RetrieveInput
	var output RetrieveOutput

	input.DataDigest = dataDigest

	for _, ip := range NodeIPs {

		worker, err := rpc.Dial("tcp", ip)
		if err != nil {
			process.GetLogger().PrintErrorTaskMessage(RetrieveTaskName, err.Error())
			continue
		}

		err = worker.Call("Retrieve.Execute", &input, &output)
		utility.CheckError(worker.Close())
		if err == nil {
			return output.RawData
		} else {
			process.GetLogger().PrintErrorTaskMessage(RetrieveTaskName, err.Error())
		}
	}

	process.GetLogger().PrintPanicErrorTaskMessage(RetrieveTaskName, "Task failed! Aborting...")

	return nil
}
