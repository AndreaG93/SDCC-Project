package Task

import (
	"SDCC-Project/BFTMapReduce/Input"
	"SDCC-Project/BFTMapReduce/Registry/WorkerMutex"
	"SDCC-Project/BFTMapReduce/Registry/WorkerResultsRegister"
)

type MapReduce struct {
}

type MapReduceInput struct {
	InputData Input.MiddleInput
}

type MapReduceOutput struct {
	Digest string
}

func (x *MapReduce) Execute(input MapReduceInput, output *MapReduceOutput) error {

	digest, rawData, err := input.InputData.PerformTask()
	if err != nil {
		return err
	}

	WorkerMutex.GetInstance().Lock()
	WorkerResultsRegister.GetInstance().Set(digest, rawData)
	WorkerMutex.GetInstance().Unlock()

	output.Digest = digest

	return nil
}
