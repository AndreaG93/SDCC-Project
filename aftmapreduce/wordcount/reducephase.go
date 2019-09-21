package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
)

func startReducePhase(guid string, AFTMapTaskOutput []*AFTMapTaskOutput) ([]*AFTReduceTaskOutput, error) {

	localityAwareReduceTaskSchedule := getLocalityAwareReduceTaskSchedule(AFTMapTaskOutput)
	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("The Locality-Aware AFT-Reduce Task Schedule for Request %s is :: %d", guid, localityAwareReduceTaskSchedule))

	startShuffleTask(AFTMapTaskOutput, localityAwareReduceTaskSchedule)

	output := reduceTask(AFTMapTaskOutput, localityAwareReduceTaskSchedule)

	return output, (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, reducePhaseComplete, utility.Encode(output))
}

func getLocalityAwareReduceTaskSchedule(input []*AFTMapTaskOutput) map[int]int {

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
