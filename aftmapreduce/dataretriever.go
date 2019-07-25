package aftmapreduce

import "SDCC-Project/aftmapreduce/registries/WorkerResultsRegister"

type DataRetriever struct {
}

type DataRetrieverInput struct {
	Digest string
}

type DataRetrieverOutput struct {
	Data []byte
}

func (x *DataRetriever) Execute(input DataRetrieverInput, output *DataRetrieverOutput) error {

	output.Data = WorkerResultsRegister.GetInstance().Get(input.Digest)
	return nil
}
