package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"sync"
)

type Request struct {
}

type RequestInput struct {
	SourceFileDigest string
}

type RequestOutput struct {
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	go manageRequest(input.SourceFileDigest)
	return nil
}

func manageRequest(digest string) {

	node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())

	splits := getSplits(digest, node.GetPropertyAsInteger(property.MapCardinality))
	arbitraryFaultTolerantMapTaskOutput := make([]*MapTaskOutput, len(splits))

	var mapWaitGroup sync.WaitGroup

	for groupId, split := range splits {

		mapWaitGroup.Add(1)

		go func(mySplit string, myGroupId int, mapWaitGroup *sync.WaitGroup) {

			arbitraryFaultTolerantMapTaskOutput[myGroupId] = NewMapTask(mySplit, myGroupId).Execute()
			mapWaitGroup.Done()

		}(split, groupId, &mapWaitGroup)
	}

	mapWaitGroup.Wait()

	mapping := getReduceTaskIndexMappedToBestGroupIndex(arbitraryFaultTolerantMapTaskOutput)

	startLocalityAwareShuffling(arbitraryFaultTolerantMapTaskOutput, mapping)
}

func getReduceTaskIndexMappedToBestGroupIndex(mapTaskOutput []*MapTaskOutput) map[int]int {

	output := make(map[int]int)

	for reduceTaskIndex := 0; reduceTaskIndex < len(mapTaskOutput); reduceTaskIndex++ {

		maxDataSize := 0

		for _, reply := range mapTaskOutput {

			currentDataSize := (*reply).MappedDataSizes[reduceTaskIndex]

			if currentDataSize > maxDataSize {
				maxDataSize = currentDataSize
				output[reduceTaskIndex] = (*reply).IdGroup
			}
		}
	}

	return output
}

func startLocalityAwareShuffling(mapTaskOutput []*MapTaskOutput, reduceTaskMappedToBestGroupId map[int]int) {

	for reduceTaskIndex, bestGroupId := range reduceTaskMappedToBestGroupId {

		receiverNodeId := (*mapTaskOutput[bestGroupId]).NodeIdsWithCorrectResult

		for _, mapOutput := range mapTaskOutput {

			if (*mapOutput).IdGroup != bestGroupId {
				sendDataTask((*mapOutput).NodeIdsWithCorrectResult, (*mapOutput).IdGroup, receiverNodeId, bestGroupId, (*mapOutput).ReplayDigest, reduceTaskIndex)
			}
		}
	}
}
