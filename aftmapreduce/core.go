package aftmapreduce

import (
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/registries/zookeeperclient"
	"SDCC-Project/utility"
	"fmt"
	"net/rpc"
	"strconv"
	"strings"
)

func ManageClientRequest(request *Request) {

	var transientData [][]byte

	faultToleranceLevel := 2
	clientData := request.getClientData()

	for {

		internetAddressesOfAvailableWorkers := zookeeperclient.GetInstance().GetMembersInternetAddress()

		switch request.getStatus() {
		case InitialPhase:

			splits, err := (*clientData).Split()
			utility.CheckError(err)

			transientData = performCurrentTask(splits, faultToleranceLevel, internetAddressesOfAvailableWorkers)
			request.Checkpoint(utility.MatrixToArray(transientData))
			continue

		case AfterMapPhase:

			splits := (*clientData).Shuffle(transientData)

			transientData := performCurrentTask(splits, faultToleranceLevel, internetAddressesOfAvailableWorkers)
			request.Checkpoint(utility.MatrixToArray(transientData))
			continue

		case AfterReducePhase:
			finalOutput := (*clientData).CollectResults(transientData)
			request.Checkpoint(finalOutput)
			return
		}
	}
}

func performCurrentTask(splits []data.TransientData, faultToleranceLevel int, workersInternetAddress map[int]string) [][]byte {

	mapReduceInternetAddresses := make([]string, len(workersInternetAddress))

	index := 0
	for key, value := range workersInternetAddress {
		mapReduceInternetAddresses[index] = fmt.Sprintf("%s:%d", value, key+MapReduceRPCBasePort)
		index++
	}

	output := make([][]byte, len(splits))

	for index := range splits {

		task := NewTask(splits[index], faultToleranceLevel, mapReduceInternetAddresses)
		digest, workerAddresses := task.Execute()

		for _, address := range workerAddresses {

			index, _ := strconv.Atoi(strings.Split(address, ":")[1])
			index = index - MapReduceRPCBasePort

			address = fmt.Sprintf("%s:%d", strings.Split(address, ":")[0], index+MapReduceGetRPCBasePort)
		}

		output[index] = retrieveDataFromWorker(digest, workerAddresses)
	}

	return output
}

func retrieveDataFromWorker(digest string, workersAddresses []string) []byte {

	for _, address := range workersAddresses {

		var input DataRetrieverInput
		var output DataRetrieverOutput

		input.Digest = digest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			continue
		}

		defer func() {
			utility.CheckError(worker.Close())
		}()

		err = worker.Call("DataRetriever.Execute", &input, &output)
		if err == nil {
			return output.Data
		}
	}
	return nil
}
