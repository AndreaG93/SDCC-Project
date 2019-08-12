package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/utility"
	"net/rpc"
)

type Retrieve struct {
}

type RetrieveInput struct {
	dataDigest string
}

type RetrieveOutput struct {
	rawData []byte
}

func (x *Retrieve) Execute(input RetrieveInput, output *RetrieveOutput) error {

	output.rawData = node.GetDataRegistry().Get(input.dataDigest).([]byte)
	return nil
}

func retrieveFrom(NodeIPs []string, dataDigest string) []byte {

	for _, ip := range NodeIPs {

		var input RetrieveInput
		var output RetrieveOutput

		input.dataDigest = dataDigest

		worker, err := rpc.Dial("tcp", ip)
		if err != nil {
			continue
		}

		err = worker.Call("DataRetriever.Execute", &input, &output)
		utility.CheckError(worker.Close())
		if err == nil {
			return output.rawData
		}
	}
	return nil
}
