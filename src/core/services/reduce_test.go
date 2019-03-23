package services

import (
	"SDCC-Project-WorkerNode/src/core/data-structuresttt"
	"os"
	"testing"
)

func Test_ReduceService(t *testing.T) {

	var data []*data_structuresttt.WordTokenList

	data = make([]*data_structuresttt.WordTokenList, 3)
	data[0] = data_structuresttt.BuildWordTokenList()
	data[1] = data_structuresttt.BuildWordTokenList()
	data[2] = data_structuresttt.BuildWordTokenList()

	(*data[0]).InsertWordToken(data_structuresttt.BuildWordToken("Andrea", 5))
	(*data[1]).InsertWordToken(data_structuresttt.BuildWordToken("Andrea", 5))
	(*data[2]).InsertWordToken(data_structuresttt.BuildWordToken("Andrea", 5))

	(*data[0]).InsertWordToken(data_structuresttt.BuildWordToken("Graziani", 5))
	(*data[1]).InsertWordToken(data_structuresttt.BuildWordToken("Graziani", 5))
	(*data[2]).InsertWordToken(data_structuresttt.BuildWordToken("Graziani", 5))

	(*data[0]).InsertWordToken(data_structuresttt.BuildWordToken("Akko", 5))
	(*data[1]).InsertWordToken(data_structuresttt.BuildWordToken("Akko", 5))
	(*data[2]).InsertWordToken(data_structuresttt.BuildWordToken("Akko", 5))

	input := ReduceInput{data}
	output := ReduceOutput{}
	obj := Reduce{}

	if myError := obj.Execute(input, &output); myError != nil {
		os.Exit(1)
	}
}
