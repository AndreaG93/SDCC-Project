/*
========================================================================================================================
Name        : core/services/reduce.go
Author      : Andrea Graziani
Description : This file includes RPC-function used to perform "Reduce" services.
========================================================================================================================
*/
package services

import (
	"SDCC-Project-WorkerNode/src/core/data-structures"
	"SDCC-Project-WorkerNode/src/core/utility"
)

type Reduce struct{}

type ReduceInput struct {
	InputFileNameString string
}

type ReduceOutput struct {
	OutputFileDigest string
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	var err error
	var inputData data_structures.WordTokenListArray
	var outputDataStructure *data_structures.WordTokenList
	var outputDataStructureSerialized data_structures.WordTokenListSerialized
	var outputDataStructureDigest string

	if inputData, err = data_structures.ReadWordTokenListArrayFromLocalFile(input.InputFileNameString); err != nil {
		return nil
	}

	outputDataStructure = data_structures.MergeAnArrayOfWordTokenLists(inputData)
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
