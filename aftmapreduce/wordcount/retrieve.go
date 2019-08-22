package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/utility"
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

	node.GetLogger().PrintInfoTaskMessage(RetrieveTaskName, fmt.Sprintf("Received a 'RETRIEVE' request -- Data digest requested is %s", input.DataDigest))

	output.RawData = node.GetDataRegistry().Get(input.DataDigest).([]byte)
	if output.RawData == nil {
		return errors.New("no data with given digest")
	}

	return nil
}

func retrieveFrom(NodeIPs []string, dataDigest string) []byte {

	node.GetLogger().PrintInfoTaskMessage(RetrieveTaskName, fmt.Sprintf("Target Nodes are %s", NodeIPs))

	var input RetrieveInput
	var output RetrieveOutput

	input.DataDigest = dataDigest

	for _, ip := range NodeIPs {

		worker, err := rpc.Dial("tcp", ip)
		if err != nil {
			node.GetLogger().PrintErrorTaskMessage(RetrieveTaskName, err.Error())
			continue
		}

		err = worker.Call("Retrieve.Execute", &input, &output)
		utility.CheckError(worker.Close())
		if err == nil {
			return output.RawData
		} else {
			node.GetLogger().PrintErrorTaskMessage(RetrieveTaskName, err.Error())
		}
	}

	node.GetLogger().PrintPanicErrorTaskMessage(RetrieveTaskName, "Task failed! Aborting...")

	return nil
}
