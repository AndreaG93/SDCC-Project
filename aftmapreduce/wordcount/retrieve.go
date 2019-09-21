package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"errors"
	"fmt"
)

type Retrieve struct {
}

type RetrieveInput struct {
	DataDigest string
}

type RetrieveOutput struct {
	RawData []byte
}

func (x *Retrieve) Execute(input RetrieveInput, output *RetrieveOutput) error {

	process.GetLogger().PrintInfoLevelLabeledMessage(RetrieveTaskName, fmt.Sprintf("Received a 'RETRIEVE' request -- Data digest requested is %s", input.DataDigest))

	output.RawData = process.GetDataRegistry().Get(input.DataDigest)
	if output.RawData == nil {
		return errors.New("no data with given digest")
	}

	return nil
}
