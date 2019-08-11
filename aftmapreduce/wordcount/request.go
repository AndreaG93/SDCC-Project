package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"fmt"
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
	fmt.Println(arbitraryFaultTolerantMapTaskOutput)
}

func startLocalityAwareShuffling(mapTaskOutput []*MapTaskOutput) {

	var bestGroupId int
	var maxDataSize int

	reduceTaskAmount := len(mapTaskOutput)

	for index := 0; index < reduceTaskAmount; index++ {

		maxDataSize = 0
		bestGroupId = 0

		for _, reply := range mapTaskOutput {

			currentDataSize := (*reply).MappedDataSizes[index]

			if currentDataSize > maxDataSize {
				maxDataSize = currentDataSize
				bestGroupId = (*reply).IdGroup
			}
		}

		for _, reply := range mapTaskOutput {

			if (*reply).IdGroup == bestGroupId {
				continue
			} else {

			}

		}
	}
}
