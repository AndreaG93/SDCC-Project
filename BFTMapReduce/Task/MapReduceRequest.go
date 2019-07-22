package Task

import (
	"SDCC-Project/BFTMapReduce"
	"SDCC-Project/BFTMapReduce/Input"
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/utility"
	"fmt"
	"net/rpc"
	"strconv"
	"strings"
)

type MapReduceRequest struct {
}

type MapReduceRequestInput struct {
	InputData Input.ApplicationInput
}

type MapReduceRequestOutput struct {
	Digest string
}

func (x *MapReduceRequest) Execute(input MapReduceRequestInput, output *MapReduceRequestOutput) error {

	zookeeperClient := zookeeper.New([]string{"localhost:2181"})
	internetAddressTable := zookeeperClient.GetMembersInternetAddress()
	zookeeperClient.CloseConnection()

	faultToleranceLevel := 2

	splits, err := input.InputData.Split()
	if err != nil {
		return err
	}

	mapTaskOutput := performCurrentTask(splits, faultToleranceLevel, internetAddressTable)

	splits = input.InputData.Shuffle(mapTaskOutput)

	reduceTaskOutput := performCurrentTask(splits, faultToleranceLevel, internetAddressTable)

	input.InputData.CollectResults(reduceTaskOutput)

	return nil
}

func performCurrentTask(splits []Input.MiddleInput, faultToleranceLevel int, workersInternetAddress map[int]string) [][]byte {

	mapReduceInternetAddresses := make([]string, len(workersInternetAddress))

	index := 0
	for key, value := range workersInternetAddress {
		mapReduceInternetAddresses[index] = fmt.Sprintf("%s:%d", value, key+BFTMapReduce.MapReduceRPCBasePort)
		index++
	}

	output := make([][]byte, len(splits))

	for index := range splits {

		task := NewBFTMapReduce(splits[index], faultToleranceLevel, mapReduceInternetAddresses)
		digest, workerAddresses := task.Execute()

		for _, address := range workerAddresses {

			index, _ := strconv.Atoi(strings.Split(address, ":")[1])
			index = index - BFTMapReduce.MapReduceRPCBasePort

			address = fmt.Sprintf("%s:%d", strings.Split(address, ":")[0], index+BFTMapReduce.MapReduceGetRPCBasePort)
		}

		output[index] = retrieveDataFromWorker(digest, workerAddresses)
	}

	return output
}

func retrieveDataFromWorker(digest string, workersAddresses []string) []byte {

	for _, address := range workersAddresses {

		var input MapReduceGetInput
		var output MapReduceGetOutput

		input.Digest = digest

		worker, err := rpc.Dial("tcp", address)
		if err != nil {
			continue
		}

		defer func() {
			utility.CheckError(worker.Close())
		}()

		err = worker.Call("MapReduceGet.Execute", &input, &output)
		if err == nil {
			return output.Data
		}
	}
	return nil
}
