package services

import (
	"core/data-structures"
	"core/utility"
	"strings"
)

type Map struct{}

type MapInput struct {
	InputString                  string
	OutputWordHashTableArraySize uint
}

type MapOutput struct {
	OutputFileDigest string
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	var outputData *data_structures.WordHashTableWithCollisionList
	var outputDataDigest string
	var err error

	outputData = data_structures.BuildWordHashTableWithCollisionList(input.OutputWordHashTableArraySize)

	wordScanner := utility.BuildWordScannerFromString(input.InputString)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = outputData.Insert(currentWord); err != nil {
			return err
		}
	}

	if outputDataDigest, err = outputData.GetDigest(); err != nil {
		return err
	}

	outputData.WriteOnLocalDisk()

	output.OutputFileDigest = outputDataDigest

	return nil
}
