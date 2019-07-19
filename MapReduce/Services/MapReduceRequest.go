package wordcount

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/MapReduce/wordcount/wordcountfile"
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

	workerAddress := []string{"127.0.0.1:10000", "127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003", "127.0.0.1:10004"}

	splits := input.inputData.Split()

	for split := range splits {

	}

	return nil
}
