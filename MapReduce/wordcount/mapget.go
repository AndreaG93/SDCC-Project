package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/worker/workermapregister"
)

type MapGet struct {
}

type MapGetInput struct {
	digest string
}

type MapGetOutput struct {
	data []byte
}

func (x *MapGet) Execute(input MapGetInput, output *MapGetOutput) error {

	data := workermapregister.GetInstance().Get(input.digest)
	output.data = data

	return nil
}
