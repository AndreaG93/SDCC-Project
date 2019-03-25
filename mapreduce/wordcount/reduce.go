package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenlist"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenlistgroup"
	"SDCC-Project-WorkerNode/utility"
	"io/ioutil"
)

type Reduce struct{}

type ReduceInput struct {
	InputFileNameString string
}

type ReduceOutput struct {
	OutputFileDigest string
}

func (x *Reduce) Execute(reduceInput ReduceInput, reduceOutput *ReduceOutput) error {

	var err error
	var rawInput []byte
	var input *wordtokenlistgroup.WordTokenListGroup
	var output *wordtokenlist.WordTokenList
	var outputSerialized []byte
	var outputDigest string

	if rawInput, err = ioutil.ReadFile(reduceInput.InputFileNameString); err != nil {
		return err
	}

	if input, err = wordtokenlistgroup.Deserialize(rawInput); err != nil {
		return err
	}

	output = input.Merge()

	if outputSerialized, err = output.Serialize(); err != nil {
		return err
	}
	if outputDigest, err = utility.GenerateDigestUsingSHA512(outputSerialized); err != nil {
		return err
	}
	if err = utility.WriteToLocalDisk(outputDigest, outputSerialized); err != nil {
		return err
	}

	reduceOutput.OutputFileDigest = outputDigest

	return nil
}
