package services

import (
	"SDCC-Project-WorkerNode/src/core/data-structures"
	"SDCC-Project-WorkerNode/src/core/utility"
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

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	var err error
	var rawInputData []byte
	var inputData string
	var outputDataStructure *data_structures.WordTokenHashTable
	var outputDataStructureSerialized data_structures.WordTokenHashTableSerialized
	var outputDataStructureDigest string

	if rawInputData, err = utility.ReadLocalFile(input.InputFileNameString); err != nil {
		return err
	}
	inputData = string(rawInputData)

	outputDataStructure = data_structures.BuildWordTokenHashTable(input.OutputWordHashTableArraySize)

	wordScanner := utility.BuildWordScannerFromString(inputData)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = outputDataStructure.InsertWord(currentWord); err != nil {
			return err
		}
	}

	outputDataStructureSerialized = outputDataStructure.Serialize()

	if outputDataStructureDigest, err = utility.SHA512(outputDataStructureSerialized); err != nil {
		return err
	}
	if err = utility.WriteToLocalDisk(outputDataStructureDigest, outputDataStructureSerialized); err != nil {
		return err
	}

	output.OutputFileDigest = outputDataStructureDigest

	return nil
}
