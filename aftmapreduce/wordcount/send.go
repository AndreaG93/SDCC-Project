package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
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

	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker IP: %s", input.ReceiversInternetAddresses))
	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source Digest: %s Destination Digest: %s ReduceTaskIndex: %d", input.SourceDataDigest, input.ReceiverAssociatedDataDigest, input.WordTokenListIndex))

	currentWordTokenHashTable := WordTokenHashTable.Deserialize(process.GetDataRegistry().Get(input.SourceDataDigest))

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
			process.GetLogger().PrintErrorTaskMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
			continue
		}

		err = worker.Call("Receive.Execute", &input, &output)
		if err != nil {
			process.GetLogger().PrintErrorTaskMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
		}
	}
}

func sendDataTask(sourceNodeIds []int, sourceGroupId int, receiverNodeIds []int, receiverGroupId int, senderDataDigest string, receiverAssociatedDataDigest string, receiverReduceTaskId int) {

	senderInternetAddresses, _ := process.GetMembershipRegister().GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(sourceGroupId, sourceNodeIds, aftmapreduce.WordCountSendRPCBasePort)
	receiverInternetAddresses, _ := process.GetMembershipRegister().GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(receiverGroupId, receiverNodeIds, aftmapreduce.WordCountReceiveRPCBasePort)

	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source worker IDs:           %d", sourceNodeIds))
	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Source worker group IDs:     %d", sourceGroupId))
	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker IDs:      %d", receiverNodeIds))
	process.GetLogger().PrintInfoTaskMessage(SendTaskName, fmt.Sprintf("Destination worker group ID: %d", receiverGroupId))

	for _, sender := range senderInternetAddresses {

		input := new(SendInput)
		output := new(SendOutput)

		(*input).SourceDataDigest = senderDataDigest
		(*input).ReceiverAssociatedDataDigest = receiverAssociatedDataDigest
		(*input).WordTokenListIndex = receiverReduceTaskId
		(*input).ReceiversInternetAddresses = receiverInternetAddresses

		worker, err := rpc.Dial("tcp", sender)
		if err != nil {
			process.GetLogger().PrintErrorTaskMessage(SendTaskName, "Send operation failed!...")
			continue
		}

		err = worker.Call("Send.Execute", input, output)
		utility.CheckError(worker.Close())
		if err == nil {
			return
		}
	}
}
