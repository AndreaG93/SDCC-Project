package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/DataStructures/wordtokenlistgroup"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/Register/Worker/WorkerReduceRegister"
	"SDCC-Project-WorkerNode/utility"
)

type Reduce struct{}

type ReduceInput struct {
	data []byte
}

type ReduceOutput struct {
	digest string
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	var err error

	tokenListGroup, err := wordtokenlistgroup.Deserialize(input.data)
	if err != nil {
		return err
	}

	data := tokenListGroup.Merge()

	rawData, err := data.Serialize()

	digest := utility.GenerateDigestUsingSHA512(rawData)

	WorkerReduceRegister.GetInstance().Set(digest, rawData)

	output.digest = digest

	return nil
}
