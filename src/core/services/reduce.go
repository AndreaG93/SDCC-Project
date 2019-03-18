/*
========================================================================================================================
Name        : core/services/reduce.go
Author      : Andrea Graziani
Description : This file includes RPC-function used to perform "Reduce" services.
========================================================================================================================
*/
package services

/*
import "core/data-structures"

type Reduce struct{}

// This structure represent the input of "Reduce" services.
type ReduceInput struct {
	Data    data_structures.WordHashTableArray
}

// This structure represent the output of "Reduce" services.
type ReduceOutput struct {
	outputFileDigest string
}

// Following function represents the published RPC routine used to perform "Reduce" services.
func (x *Reduce) Execute(pInput ReduceInput, pOutput *ReduceOutput) error {

	output := make(map[string]uint)

	for j := 0; j < len(pInput.Data); j++ {

		for word, occurrences := range pInput.Data[j] {

			output[word] += occurrences

		}
	}



	return nil
}
*/
