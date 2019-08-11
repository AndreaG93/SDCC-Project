package aftmapreduce

/*
import (
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount"
	"SDCC-Project/cloud/amazons3"
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

	transientData = nil

	for {

		internetAddressesOfAvailableWorkers := node.GetZookeeperClient().GetMembersInternetAddress()

		switch request.getStatus() {
		case InitialPhase:

			node.GetLogger().PrintMessage(fmt.Sprintf("Enter into 'InitialPhase' for request ID: %s", request.digest))

			splits, err := (*clientData).Split()
			utility.CheckError(err)

			transientData = performCurrentTask(splits, faultToleranceLevel, internetAddressesOfAvailableWorkers)
			request.Checkpoint(utility.MatrixToArray(transientData))
			continue

		case AfterMapPhase:

			amazonS3Client := amazons3.New()
			amazonS3Client.Delete()

			node.GetLogger().PrintMessage(fmt.Sprintf("Enter into 'AfterMapPhase' for request ID: %s", request.digest))

			if transientData == nil {
				transientData = utility.ArrayToMatrix(request.GetDataFromCheckpoint())
			}

			splits := (*clientData).Shuffle(transientData)

			transientData = performCurrentTask(splits, faultToleranceLevel, internetAddressesOfAvailableWorkers)
			request.Checkpoint(utility.MatrixToArray(transientData))

			continue

		case AfterReducePhase:

			node.GetLogger().PrintMessage(fmt.Sprintf("Enter into 'AfterReducePhase' for request ID: %s", request.digest))

			if transientData == nil {
				transientData = utility.ArrayToMatrix(request.GetDataFromCheckpoint())
			}

			finalOutput := (*clientData).CollectResults(transientData)
			request.Checkpoint(finalOutput)

			node.GetLogger().PrintMessage(fmt.Sprintf("Request ID %s satisfied", request.digest))
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

		var input wordcount.DataRetrieverInput
		var output wordcount.DataRetrieverOutput

		input.Digest = digest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			continue
		}

		err = worker.Call("DataRetriever.Execute", &input, &output)
		utility.CheckError(worker.Close())
		if err == nil {

			return output.Data
		}
	}
	return nil
}
*/
