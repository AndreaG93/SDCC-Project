package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"errors"
	"fmt"
	"net/rpc"
)

func startCollectPhase(guid string, reducePhaseOutput []*AFTReduceTaskOutput) error {

	var err error

	if dataArray, err := collectReduceTaskOutputFromWPG(reducePhaseOutput); err == nil {

		output := computeFinalOutputTask(dataArray)
		outputSerialized := output.Serialize()
		outputGUID := utility.GenerateDigestUsingSHA512(outputSerialized)

		if err = (*process.GetStorageKeyValueRegister()).Put(outputGUID, outputSerialized); err == nil {
			if err = (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, collectPhaseComplete, nil); err == nil {
				err = (*process.GetSystemCoordinator()).RegisterClientRequestAsComplete(guid, outputGUID)
			}
		}
	}

	return err
}

func collectReduceTaskOutputFromWPG(reducePhaseOutput []*AFTReduceTaskOutput) ([]*WordTokenList.WordTokenList, error) {

	output := make([]*WordTokenList.WordTokenList, len(reducePhaseOutput))

	for index, aftReduceTaskOutput := range reducePhaseOutput {

		targetNodeIP, _ := (*process.GetMembershipRegister()).GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(aftReduceTaskOutput.IdGroup, aftReduceTaskOutput.NodeIdsWithCorrectResult, aftmapreduce.WordCountRetrieverRPCBasePort)

		if rawData, err := retrieveFrom(targetNodeIP, aftReduceTaskOutput.ReplayDigest); err != nil {
			return nil, err
		} else {
			serializedData := WordTokenList.Deserialize(rawData)
			output[index] = serializedData
		}
	}

	return output, nil
}

func retrieveFrom(NodeIPs []string, digest string) ([]byte, error) {

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Attempting to retrieve data with digest %s from %s", digest, NodeIPs))

	var input RetrieveInput
	var output RetrieveOutput

	input.DataDigest = digest

	for _, ip := range NodeIPs {

		if worker, err := rpc.Dial("tcp", ip); err != nil {
			process.GetLogger().PrintErrorLevelMessage(err.Error())
		} else {

			if err = worker.Call("Retrieve.Execute", &input, &output); err != nil {
				process.GetLogger().PrintErrorLevelMessage(err.Error())
			} else {

				if err = worker.Close(); err == nil {
					return output.RawData, nil
				}
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("FAILED during retrieving data with digest %s from %s", digest, NodeIPs))
}
