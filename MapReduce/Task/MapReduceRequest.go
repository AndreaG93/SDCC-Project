package Task

import (
	"SDCC-Project/MapReduce/Input"
	"SDCC-Project/utility"
	"net/rpc"
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

	workerAddress := []string{"127.0.0.1:12001", "127.0.0.1:12002", "127.0.0.1:12003", "127.0.0.1:12004", "127.0.0.1:12005", "127.0.0.1:12006"}
	faultToleranceLevel := 2

	splits, err := input.InputData.Split()
	if err != nil {
		return err
	}

	mapTaskOutput := performCurrentTask(splits, faultToleranceLevel, workerAddress)

	splits = input.InputData.Shuffle(mapTaskOutput)

	reduceTaskOutput := performCurrentTask(splits, faultToleranceLevel, workerAddress)

	input.InputData.CollectResults(reduceTaskOutput)

	return nil
}

func performCurrentTask(splits []Input.MiddleInput, faultToleranceLevel int, workerAddress []string) [][]byte {

	//var myWaitGroup sync.WaitGroup
	output := make([][]byte, len(splits))

	for index := range splits {

		task := NewBFTMapReduce(splits[index], faultToleranceLevel, workerAddress)
		digest, workerAddresses := task.Execute()

		output[index] = retrieveDataFromWorker(digest, workerAddresses)

		//myWaitGroup.Add(1)
		/*
			go func(x int) {
				task := NewBFTMapReduce(splits[index], faultToleranceLevel, workerAddress)
				digest, workerAddresses := task.Execute()

				output[x] = retrieveDataFromWorker(digest, workerAddresses)

				myWaitGroup.Done()
			}(index)
		*/

	}

	//myWaitGroup.Wait()
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
