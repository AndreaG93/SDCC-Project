/*
========================================================================================================================
Name        : core/services/reduce.go
Author      : Andrea Graziani
Description : This file includes RPC-function used to perform "Reduce" services.
========================================================================================================================
*/
package services

import (
	"core/data-structures"
	"core/utility"
)

type Reduce struct{}

// This structure represent the input of "Reduce" services.
type ReduceInput struct {
	Data []*data_structures.WordTokenList
}

// This structure represent the output of "Reduce" services.
type ReduceOutput struct {
	Digest string
}

// Following function represents the published RPC routine used to perform "Reduce" services.
func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	var digest string
	var err error
	var currentWordTokenList *data_structures.WordTokenList
	var outputWordTokenList *data_structures.WordTokenList

	outputWordTokenList = data_structures.BuildWordTokenList()

	for index := 0; index < len(input.Data); index++ {

		currentWordTokenList = input.Data[index]

		(*outputWordTokenList).Merge(currentWordTokenList)
	}

	(outputWordTokenList).Print()

	if digest, err = utility.SHA512((*outputWordTokenList).Serialize()); err != nil {
		return err
	}

	(*output).Digest = digest

	return nil
}
