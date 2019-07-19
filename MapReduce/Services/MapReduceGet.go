package wordcount

import "SDCC-Project/MapReduce/Registry/WorkerResultsRegister"

type MapReduceGet struct {
}

type MapReduceGetInput struct {
	Digest string
}

type MapReduceGetOutput struct {
	Data []byte
}

func (x *MapReduceGet) Execute(input MapReduceGetInput, output *MapReduceGetOutput) error {

	output.Data = WorkerResultsRegister.GetInstance().Get(input.Digest)

	return nil
}
