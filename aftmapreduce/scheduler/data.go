package scheduler

import (
	"SDCC-Project/aftmapreduce/worker"
	"sync"
)

func StartLocalityAwareShuffleTask(splitCardinality int, workerMapResponses []*worker.WorkerMapResponse) {

	var waitGroup sync.WaitGroup

	for index := 0; index < splitCardinality; index++ {
		waitGroup.Add(1)
		go shuffle(index, workerMapResponses, &waitGroup)
	}

	waitGroup.Wait()
}

func shuffle(splitIndex int, workerMapResponses []*worker.WorkerMapResponse, waitGroup *sync.WaitGroup) {

	maxDataSize := 0
	candidateNodeId := -1

	for index := 0; index < len(workerMapResponses); index++ {

		currentDataSize := (*workerMapResponses[index]).DataSizes[splitIndex]

		if currentDataSize > maxDataSize {
			maxDataSize = currentDataSize
			candidateNodeId = index
		}
	}

	waitGroup.Done()
}

type WorkerReduceResponse struct {
	workerInternetAddress string
	DataSizes             map[uint]uint
}
