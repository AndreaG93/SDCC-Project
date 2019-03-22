package services

import (
	"SDCC-Project-WorkerNode/src/core/data-structures"
	"os"
	"testing"
)

func Test_ReduceService(t *testing.T) {

	var data []*data_structures.WordTokenList

	data = make([]*data_structures.WordTokenList, 3)
	data[0] = data_structures.BuildWordTokenList()
	data[1] = data_structures.BuildWordTokenList()
	data[2] = data_structures.BuildWordTokenList()

	(*data[0]).InsertWordToken(data_structures.BuildWordToken("Andrea", 5))
	(*data[1]).InsertWordToken(data_structures.BuildWordToken("Andrea", 5))
	(*data[2]).InsertWordToken(data_structures.BuildWordToken("Andrea", 5))

	(*data[0]).InsertWordToken(data_structures.BuildWordToken("Graziani", 5))
	(*data[1]).InsertWordToken(data_structures.BuildWordToken("Graziani", 5))
	(*data[2]).InsertWordToken(data_structures.BuildWordToken("Graziani", 5))

	(*data[0]).InsertWordToken(data_structures.BuildWordToken("Akko", 5))
	(*data[1]).InsertWordToken(data_structures.BuildWordToken("Akko", 5))
	(*data[2]).InsertWordToken(data_structures.BuildWordToken("Akko", 5))

	input := ReduceInput{data}
	output := ReduceOutput{}
	obj := Reduce{}

	if myError := obj.Execute(input, &output); myError != nil {
		os.Exit(1)
	}
}
