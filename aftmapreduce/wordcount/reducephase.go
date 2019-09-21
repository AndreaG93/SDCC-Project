package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
)

func startReducePhase(guid string, AFTMapTaskOutput []*AFTMapTaskOutput) ([]*AFTReduceTaskOutput, error) {

	localityAwarenessData := getLocalityAwareReduceTaskMappedToNodeGroupId(AFTMapTaskOutput)
	localityAwareShuffleTask(AFTMapTaskOutput, localityAwarenessData)

	output := reduceTask(AFTMapTaskOutput, localityAwarenessData)

	return output, (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, reducePhaseComplete, utility.Encode(output))
}

func startCollectPhase(guid string, AFTReduceTaskOutput []*AFTReduceTaskOutput) error {

	var err error

	dataArray := retrieveTask(AFTReduceTaskOutput)

	output := computeFinalOutputTask(dataArray)
	outputSerialized := output.Serialize()
	outputGUID := utility.GenerateDigestUsingSHA512(outputSerialized)

	if err = (*process.GetStorageKeyValueRegister()).Put(outputGUID, outputSerialized); err == nil {
		if err = (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, collectPhaseComplete, nil); err == nil {
			err = (*process.GetSystemCoordinator()).RegisterClientRequestAsComplete(guid, outputGUID)
		}
	}

	return err
}
