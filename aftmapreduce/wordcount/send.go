package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"net/rpc"
)

type Send struct {
}

type SendInput struct {
	SourceDataDigest           string
	WordTokenListIndex         int
	receiversInternetAddresses []string
}

type SendOutput struct {
	SendDataDigest string
}

func (x *Send) Execute(input SendInput, output *SendOutput) error {

	data := (node.GetCache().Get(input.SourceDataDigest)).(*WordTokenHashTable.WordTokenHashTable).GetWordTokenListAt(input.WordTokenListIndex)
	dataDigest := data.GetDigest()

	sendDataToWorker(data, dataDigest, input.receiversInternetAddresses)

	output.SendDataDigest = dataDigest

	return nil
}

func sendDataToWorker(data *WordTokenList.WordTokenList, dataDigest string, receiversInternetAddresses []string) {

	for _, address := range receiversInternetAddresses {

		var input ReceiveInput
		var output ReceiveOutput

		input.ReceivedDataDigest = dataDigest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			continue
		}

		err = worker.Call("Receive.Execute", &input, &output)
		utility.CheckError(worker.Close())
	}
}

func sendDataTask(sourceNodeIds []int, sourceGroupId int, receiverGroupId, dataDigest string) {

}
