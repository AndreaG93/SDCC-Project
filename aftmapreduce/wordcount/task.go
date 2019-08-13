package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"fmt"
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

func mapTask(input []string) []*AFTMapTaskOutput {

	output := make([]*AFTMapTaskOutput, len(input))

	var mapWaitGroup sync.WaitGroup

	for index, split := range input {

		mapWaitGroup.Add(1)
		go func(mySplit string, myGroupId int) {

			output[myGroupId] = NewMapTask(mySplit, myGroupId).Execute()
			mapWaitGroup.Done()

		}(split, index)
	}

	mapWaitGroup.Wait()
	node.GetLogger().PrintInfoCompleteTaskMessage(MapTaskName)
	return output
}

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

	for index, bestGroupId := range reduceTaskMappedToNodeGroupId {

		receiverDigestData := input[bestGroupId].ReplayDigest
		receiverNodeId := (*input[bestGroupId]).NodeIdsWithCorrectResult

		for _, mapOutput := range input {

			if (*mapOutput).IdGroup != bestGroupId {
				sendDataTask((*mapOutput).NodeIdsWithCorrectResult, (*mapOutput).IdGroup, receiverNodeId, bestGroupId, (*mapOutput).ReplayDigest, receiverDigestData, index)
			}
		}
	}
	node.GetLogger().PrintInfoCompleteTaskMessage(ShuffleTaskName)
}

func reduceTask(input []*AFTMapTaskOutput, reduceTaskMappedToNodeGroupId map[int]int) []*AFTReduceTaskOutput {

	output := make([]*AFTReduceTaskOutput, len(input))
	var mapWaitGroup sync.WaitGroup

	for index, bestGroupId := range reduceTaskMappedToNodeGroupId {

		receiverDigestData := input[bestGroupId].ReplayDigest
		receiverNodeId := (*input[bestGroupId]).NodeIdsWithCorrectResult

		mapWaitGroup.Add(1)
		go func(targetNodeIds []int, targetNodeGroupId int, reduceTaskIdentifierDigest string, reduceTaskIndex int) {

			output[reduceTaskIndex] = NewAFTReduceTask(targetNodeIds, targetNodeGroupId, reduceTaskIdentifierDigest, reduceTaskIndex).Execute()
			mapWaitGroup.Done()

		}(receiverNodeId, bestGroupId, receiverDigestData, index)
	}

	mapWaitGroup.Wait()

	node.GetLogger().PrintInfoTaskMessage(ReduceTaskName, fmt.Sprintf("Output: %s", output))
	node.GetLogger().PrintInfoCompleteTaskMessage(ReduceTaskName)
	return output
}

func retrieveTask(input []*AFTReduceTaskOutput) []*WordTokenList.WordTokenList {

	output := make([]*WordTokenList.WordTokenList, len(input))

	for index, aftReduceTaskOutput := range input {

		targetNodeIP := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(aftReduceTaskOutput.IdGroup, aftmapreduce.WordCountRetrieverRPCBasePort, aftReduceTaskOutput.NodeIdsWithCorrectResult)

		rawData := retrieveFrom(targetNodeIP, aftReduceTaskOutput.ReplayDigest)
		serializedData, err := WordTokenList.Deserialize(rawData)
		utility.CheckError(err)

		output[index] = serializedData
	}
	node.GetLogger().PrintInfoCompleteTaskMessage(RetrieveTaskName)
	return output
}

func computeFinalOutputTask(input []*WordTokenList.WordTokenList) *WordTokenList.WordTokenList {

	output := input[0]

	for index := 1; index < len(input); index++ {
		output.Merge(input[index])
	}

	return output
}
