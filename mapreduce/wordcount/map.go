package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/utility"
	"io/ioutil"
	"strings"
)

type Map struct {
}

type MapInput struct {
	InputFileNameString          string
	OutputWordHashTableArraySize uint
}

type MapOutput struct {
	OutputFileDigest string
}

func (x *Map) Execute(mapInput MapInput, mapOutput *MapOutput) error {

	var err error
	var rawInput []byte
	var output *wordtokenhashtable.WordTokenHashTable
	var outputSerialized []byte
	var outputDigest string

	if rawInput, err = ioutil.ReadFile(mapInput.InputFileNameString); err != nil {
		return err
	}

	output = wordtokenhashtable.New(mapInput.OutputWordHashTableArraySize)

	wordScanner := utility.BuildWordScannerFromString(string(rawInput))

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = output.InsertWord(currentWord); err != nil {
			return err
		}
	}

	if outputSerialized, err = output.Serialize(); err != nil {
		return err
	}
	if outputDigest, err = utility.GenerateDigestUsingSHA512(outputSerialized); err != nil {
		return err
	}
	if err = ioutil.WriteFile(outputDigest, outputSerialized, 0777); err != nil {
		return err
	}

	mapOutput.OutputFileDigest = outputDigest

	return nil
}
