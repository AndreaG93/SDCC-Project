package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"fmt"
	"net/rpc"
)

type Send struct {
}

type SendInput struct {
	SourceDataDigest             string
	ReceiverAssociatedDataDigest string
	WordTokenListIndex           int
	ReceiversInternetAddresses   []string
}

type SendOutput struct {
	SendDataDigest string
}

func (x *Send) Execute(input SendInput, output *SendOutput) error {

	node.GetLogger().PrintMessage(fmt.Sprintf("Destination worker IP: %s", input.ReceiversInternetAddresses))
	node.GetLogger().PrintMessage(fmt.Sprintf("Source Digest: %s Destination Digest: %s ReduceTaskIndex: %d", input.SourceDataDigest, input.ReceiverAssociatedDataDigest, input.WordTokenListIndex))

	data := (node.GetDataRegistry().Get(input.SourceDataDigest)).(*WordTokenHashTable.WordTokenHashTable).GetWordTokenListAt(input.WordTokenListIndex)
	dataDigest := data.GetDigest()

	sendDataToWorker(data, dataDigest, input.ReceiverAssociatedDataDigest, input.ReceiversInternetAddresses)

	output.SendDataDigest = dataDigest

	return nil
}

func sendDataToWorker(data *WordTokenList.WordTokenList, dataDigest string, receiverAssociatedDataDigest string, receiversInternetAddresses []string) {

	for _, address := range receiversInternetAddresses {

		var input ReceiveInput
		var output ReceiveOutput

		input.ReceivedDataDigest = dataDigest
		input.Data, _ = data.Serialize()
		input.AssociatedDataDigest = receiverAssociatedDataDigest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			node.GetLogger().PrintMessage(fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
			continue
		}

		err = worker.Call("Receive.Execute", &input, &output)
		if err != nil {
			node.GetLogger().PrintMessage(fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
			continue
		}
		//utility.CheckError(worker.Close())
	}
}

func sendDataTask(sourceNodeIds []int, sourceGroupId int, receiverNodeIds []int, receiverGroupId int, senderDataDigest string, receiverAssociatedDataDigest string, receiverReduceTaskId int) {

	senderInternetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(sourceGroupId, aftmapreduce.WordCountSendRPCBasePort, sourceNodeIds)
	receiverInternetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(receiverGroupId, aftmapreduce.WordCountReceiveRPCBasePort, receiverNodeIds)

	node.GetLogger().PrintMessage(fmt.Sprintf("Send Task:\\n\\tSource\\n\\t\\tGroup Node ID: %d\n\t\tNode IDs: %d", sourceGroupId, sourceNodeIds))

	node.GetLogger().PrintMessage("Destination worker IDs:" + fmt.Sprintf("%d", receiverGroupId))
	node.GetLogger().PrintMessage(fmt.Sprintf("Destination worker Group ID: %d", receiverGroupId))

	node.GetLogger().PrintMessage("Destination worker internetAddresses:" + fmt.Sprintf("%s", receiverInternetAddresses))
	node.GetLogger().PrintMessage("Source worker internetAddresses:" + fmt.Sprintf("%s", senderInternetAddresses))

	for _, sender := range senderInternetAddresses {

		input := new(SendInput)
		output := new(SendOutput)

		(*input).SourceDataDigest = senderDataDigest
		(*input).ReceiverAssociatedDataDigest = receiverAssociatedDataDigest
		(*input).WordTokenListIndex = receiverReduceTaskId
		(*input).ReceiversInternetAddresses = receiverInternetAddresses
		fmt.Println(receiverInternetAddresses)

		worker, err := rpc.Dial("tcp", sender)
		utility.CheckError(err)

		err = worker.Call("Send.Execute", input, output)
		utility.CheckError(worker.Close())
		if err == nil {
			return
		}
	}
}
