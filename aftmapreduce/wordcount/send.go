package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
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

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Received SEND order! Destination worker IPs: %s :: Source Digest: %s :: Destination Digest: %s :: ReduceTaskIndex: %d", input.ReceiversInternetAddresses, input.SourceDataDigest, input.ReceiverAssociatedDataDigest, input.WordTokenListIndex))

	rawData := process.GetDataRegistry().Get(input.SourceDataDigest)

	if input.WordTokenListIndex == -1 {
		(*output).SendDataDigest = input.SourceDataDigest
	} else {

		if localWordTokenHashTable, err := WordTokenHashTable.Deserialize(rawData); err != nil {
			return err
		} else {
			(*output).SendDataDigest, rawData, err = localWordTokenHashTable.GetWordTokenListAt(input.WordTokenListIndex).GetDigestAndSerializedData()
		}
	}

	sendDataToWorker(rawData, (*output).SendDataDigest, input.ReceiverAssociatedDataDigest, input.ReceiversInternetAddresses)

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
			process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
			continue
		}

		err = worker.Call("Receive.Execute", &input, &output)
		if err != nil {
			process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Error during data transmission to %s: %s", address, err.Error()))
		}
	}
}

func sendDataTask(sourceNodeIds []int, sourceGroupId int, receiverNodeIds []int, receiverGroupId int, senderDataDigest string, receiverAssociatedDataDigest string, receiverReduceTaskId int) {

	senderInternetAddresses, _ := process.GetMembershipRegister().GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(sourceGroupId, sourceNodeIds, aftmapreduce.WordCountSendRPCBasePort)
	receiverInternetAddresses, _ := process.GetMembershipRegister().GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(receiverGroupId, receiverNodeIds, aftmapreduce.WordCountReceiveRPCBasePort)

	process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Source worker IDs:           %d", sourceNodeIds))
	process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Source worker group IDs:     %d", sourceGroupId))
	process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Destination worker IDs:      %d", receiverNodeIds))
	process.GetLogger().PrintInfoLevelLabeledMessage(SendTaskName, fmt.Sprintf("Destination worker group ID: %d", receiverGroupId))

	for _, sender := range senderInternetAddresses {

		input := new(SendInput)
		output := new(SendOutput)

		(*input).SourceDataDigest = senderDataDigest
		(*input).ReceiverAssociatedDataDigest = receiverAssociatedDataDigest
		(*input).WordTokenListIndex = receiverReduceTaskId
		(*input).ReceiversInternetAddresses = receiverInternetAddresses

		if worker, err := rpc.Dial("tcp", sender); err != nil {
			process.GetLogger().PrintInfoLevelMessage(err.Error())
		} else {
			if err = worker.Call("Send.Execute", input, output); err != nil {
				process.GetLogger().PrintInfoLevelMessage(err.Error())
			} else {
				if err := worker.Close(); err == nil {
					return
				}
			}
		}
	}
}
