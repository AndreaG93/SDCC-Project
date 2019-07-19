package wordcount

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/MapReduce/Task/BFTMapTask"
	"SDCC-Project/utility"
	"net/rpc"
	"sync"
)

type MapReduceRequest struct {
}

type MapReduceRequestInput struct {
	inputData Data.RawInput
}

type MapReduceRequestOutput struct {
	Digest string
}

func (x *MapReduceRequest) Execute(input MapReduceRequestInput, output *MapReduceRequestOutput) error {

	var myWaitGroup sync.WaitGroup

	workerAddress := []string{"127.0.0.1:10000", "127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003", "127.0.0.1:10004"}
	faultToleranceLevel := 1

	splits := input.inputData.Split()
	mapTaskOutput := make([][]byte, len(splits))

	for index := range splits {

		myWaitGroup.Add(1)

		go func(x int) {
			task := BFTMapTask.New(splits[index], faultToleranceLevel, workerAddress)
			digest, workerAddresses := task.Execute()

			mapTaskOutput[x] = retrieveDataFromWorker(digest, workerAddresses)

			myWaitGroup.Done()
		}(index)
	}

	myWaitGroup.Wait()

	return nil
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

		err = worker.Call("Services.MapReduceGet", &input, &output)
		if err != nil {
			return output.Data
		}
	}
	return nil
}
