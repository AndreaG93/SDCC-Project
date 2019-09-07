package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
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

	var rawData []byte
	var dataDigest string

	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker IP: %s", input.ReceiversInternetAddresses))
	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source Digest: %s Destination Digest: %s ReduceTaskIndex: %d", input.SourceDataDigest, input.ReceiverAssociatedDataDigest, input.WordTokenListIndex))

	currentWordTokenHashTable := WordTokenHashTable.Deserialize(node.GetDataRegistry().Get(input.SourceDataDigest))

	if input.WordTokenListIndex == -1 {

		rawData = currentWordTokenHashTable.Serialize()
		dataDigest = input.SourceDataDigest

	} else {
		data := currentWordTokenHashTable.GetWordTokenListAt(input.WordTokenListIndex)
		dataDigest, rawData = data.GetDigestAndSerializedData()
	}

	sendDataToWorker(rawData, dataDigest, input.ReceiverAssociatedDataDigest, input.ReceiversInternetAddresses)

	output.SendDataDigest = dataDigest

	return nil
}

func sendDataToWorker(data []byte, dataDigest string, receiverAssociatedDataDigest string, receiversInternetAddresses []string) {

	for _, address := range receiversInternetAddresses {

		var input ReceiveInput
		var output ReceiveOutput

		input.ReceivedDataDigest = dataDigest
		input.Data = data
		input.AssociatedDataDigest = receiverAssociatedDataDigest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			node.GetLogger().PrintPanicErrorTaskMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
			continue
		}

		err = worker.Call("Receive.Execute", &input, &output)
		if err != nil {
			node.GetLogger().PrintPanicErrorTaskMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
		}
	}
}

func sendDataTask(sourceNodeIds []int, sourceGroupId int, receiverNodeIds []int, receiverGroupId int, senderDataDigest string, receiverAssociatedDataDigest string, receiverReduceTaskId int) {

	senderInternetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(sourceGroupId, aftmapreduce.WordCountSendRPCBasePort, sourceNodeIds)
	receiverInternetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(receiverGroupId, aftmapreduce.WordCountReceiveRPCBasePort, receiverNodeIds)

	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source worker IDs:           %d", sourceNodeIds))
	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source worker group IDs:     %d", sourceGroupId))
	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker IDs:      %d", receiverNodeIds))
	node.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker group ID: %d", receiverGroupId))

	for _, sender := range senderInternetAddresses {

		input := new(SendInput)
		output := new(SendOutput)

		(*input).SourceDataDigest = senderDataDigest
		(*input).ReceiverAssociatedDataDigest = receiverAssociatedDataDigest
		(*input).WordTokenListIndex = receiverReduceTaskId
		(*input).ReceiversInternetAddresses = receiverInternetAddresses

		worker, err := rpc.Dial("tcp", sender)
		utility.CheckError(err)

		err = worker.Call("Send.Execute", input, output)
		utility.CheckError(worker.Close())
		if err == nil {
			return
		}
	}

	node.GetLogger().PrintPanicErrorTaskMessage(SendTaskName, "Send operation failed! aborting...")
}
