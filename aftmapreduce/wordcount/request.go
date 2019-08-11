package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount/aft"
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

	splits := getSplits(digest, node.GetPropertyAsInteger(property.MapCardinality))
	arbitraryFaultTolerantMapTaskOutput := make([]*aft.MapTaskOutput, len(splits))

	var mapWaitGroup sync.WaitGroup

	for index, split := range splits {

		mapWaitGroup.Add(1)

		go func(index int, mapWaitGroup *sync.WaitGroup) {

			arbitraryFaultTolerantMapTaskOutput[index] = aft.NewMapTask(split, index).Execute()
			mapWaitGroup.Done()

		}(index, &mapWaitGroup)
	}

	mapWaitGroup.Wait()
	fmt.Println(arbitraryFaultTolerantMapTaskOutput)
}

func startLocalityAwareShuffling(mapTaskOutput []*aft.MapTaskOutput) {

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
