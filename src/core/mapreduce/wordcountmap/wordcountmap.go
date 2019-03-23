package wordcountmap

import (
	"SDCC-Project-WorkerNode/src/core/datastructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/src/core/utility"
	"io/ioutil"
	"strings"
)

type Map struct {
}

type Input struct {
	InputFileNameString          string
	OutputWordHashTableArraySize uint
}

type Output struct {
	OutputFileDigest string
}

func (x *Map) Execute(input Input, output *Output) error {

	var err error
	var rawInputData []byte
	var inputData string

	var outputDataStructure *wordtokenhashtable.WordTokenHashTable
	var outputDataStructureDigest string


	if rawInputData, err = ioutil.ReadFile(input.InputFileNameString); err != nil {
		return err
	}



	

	outputDataStructure = wordtokenhashtable.New(input.OutputWordHashTableArraySize)

	wordScanner := utility.BuildWordScannerFromString(inputData)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = outputDataStructure.InsertWord(currentWord); err != nil {
			return err
		}
	}

	outputDataStructureSerialized = outputDataStructure.Serialize()

	if outputDataStructureDigest, err = utility.GenerateDigestOfDataUsingSHA512(outputDataStructureSerialized); err != nil {
		return err
	}
	if err = utility.WriteToLocalDisk(outputDataStructureDigest, outputDataStructureSerialized); err != nil {
		return err
	}

	output.OutputFileDigest = outputDataStructureDigest

	return nil
}
