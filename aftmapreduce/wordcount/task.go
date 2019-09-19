package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"sync"
)

const (
	ReduceTaskName   = "REDUCE"
	MapTaskName      = "MAP"
	ReceiveTaskName  = "RECEIVE"
	SendTaskName     = "SEND"
	RetrieveTaskName = "RETRIEVE"
	ShuffleTaskName  = "SHUFFLE"
)

func getLocalityAwareReduceTaskMappedToNodeGroupId(input []*AFTMapTaskOutput) map[int]int {

	output := make(map[int]int)

	for reduceTaskIndex := 0; reduceTaskIndex < len(input); reduceTaskIndex++ {

		maxDataSize := 0

		for _, reply := range input {

			currentDataSize := (*reply).MappedDataSizes[reduceTaskIndex]

			if currentDataSize > maxDataSize {
				maxDataSize = currentDataSize
				output[reduceTaskIndex] = (*reply).IdGroup
			}
		}
	}

	return output
}

func localityAwareShuffleTask(input []*AFTMapTaskOutput, reduceTaskMappedToNodeGroupId map[int]int) {

	for _, mapOutput := range input {

		nodeWithoutCorrectData := make([]int, 0)

		nodeWithCorrectData := (*mapOutput).NodeIdsWithCorrectResult
		allNode := (*process.GetMembershipRegister()).GetWorkerProcessIDs(mapOutput.IdGroup)

		for _, nodeID := range allNode {

			found := false

			for _, nodeIDWithCorrectData := range nodeWithCorrectData {

				if nodeID == nodeIDWithCorrectData {
					found = true
				}
			}

			if !found {
				nodeWithoutCorrectData = append(nodeWithoutCorrectData, nodeID)
			}
		}

		sendDataTask((*mapOutput).NodeIdsWithCorrectResult, (*mapOutput).IdGroup, nodeWithoutCorrectData, (*mapOutput).IdGroup, (*mapOutput).ReplayDigest, "", -1)
	}

	for index, bestGroupId := range reduceTaskMappedToNodeGroupId {

		receiverDigestData := input[bestGroupId].ReplayDigest
		receiverNodeId := (*process.GetMembershipRegister()).GetWorkerProcessIDs(bestGroupId)

		for _, mapOutput := range input {

			if (*mapOutput).IdGroup != bestGroupId {
				sendDataTask((*mapOutput).NodeIdsWithCorrectResult, (*mapOutput).IdGroup, receiverNodeId, bestGroupId, (*mapOutput).ReplayDigest, receiverDigestData, index)
			}
		}
	}
	process.GetLogger().PrintInfoCompleteTaskMessage(ShuffleTaskName)
}

func reduceTask(input []*AFTMapTaskOutput, reduceTaskMappedToNodeGroupId map[int]int) []*AFTReduceTaskOutput {

	output := make([]*AFTReduceTaskOutput, len(input))
	var mapWaitGroup sync.WaitGroup

	for index, bestGroupId := range reduceTaskMappedToNodeGroupId {

		receiverDigestData := input[bestGroupId].ReplayDigest

		mapWaitGroup.Add(1)
		go func(targetNodeGroupId int, reduceTaskIdentifierDigest string, reduceTaskIndex int) {

			output[reduceTaskIndex] = Execute(NewAFTReduceTask(targetNodeGroupId, reduceTaskIdentifierDigest, reduceTaskIndex)).(*AFTReduceTaskOutput)
			mapWaitGroup.Done()

		}(bestGroupId, receiverDigestData, index)
	}

	mapWaitGroup.Wait()

	process.GetLogger().PrintInfoCompleteTaskMessage(ReduceTaskName)
	return output
}

func retrieveTask(input []*AFTReduceTaskOutput) []*WordTokenList.WordTokenList {

	output := make([]*WordTokenList.WordTokenList, len(input))

	for index, aftReduceTaskOutput := range input {

		targetNodeIP, _ := (*process.GetMembershipRegister()).GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(aftReduceTaskOutput.IdGroup, aftReduceTaskOutput.NodeIdsWithCorrectResult, aftmapreduce.WordCountRetrieverRPCBasePort)

		rawData := retrieveFrom(targetNodeIP, aftReduceTaskOutput.ReplayDigest)
		serializedData := WordTokenList.Deserialize(rawData)

		output[index] = serializedData
	}
	process.GetLogger().PrintInfoCompleteTaskMessage(RetrieveTaskName)
	return output
}

func computeFinalOutputTask(input []*WordTokenList.WordTokenList) *WordTokenList.WordTokenList {

	output := input[0]

	for index := 1; index < len(input); index++ {
		output.Merge(input[index])
	}

	return output
}
