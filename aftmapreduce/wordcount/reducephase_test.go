package wordcount

import (
	"fmt"
	"testing"
)

func Test_LocalityAwareShuffle1(t *testing.T) {

	input := make([]*AFTMapTaskOutput, 3)

	input[0] = &AFTMapTaskOutput{
		IdGroup:                  0,
		ReplayDigest:             "",
		NodeIdsWithCorrectResult: nil,
		MappedDataSizes:          make(map[int]int),
	}
	input[1] = &AFTMapTaskOutput{
		IdGroup:                  1,
		ReplayDigest:             "",
		NodeIdsWithCorrectResult: nil,
		MappedDataSizes:          make(map[int]int),
	}
	input[2] = &AFTMapTaskOutput{
		IdGroup:                  2,
		ReplayDigest:             "",
		NodeIdsWithCorrectResult: nil,
		MappedDataSizes:          make(map[int]int),
	}

	(*input[0]).MappedDataSizes[0] = 10
	(*input[0]).MappedDataSizes[1] = 10
	(*input[0]).MappedDataSizes[2] = 10

	(*input[1]).MappedDataSizes[0] = 60
	(*input[1]).MappedDataSizes[1] = 15
	(*input[1]).MappedDataSizes[2] = 15

	(*input[2]).MappedDataSizes[0] = 10
	(*input[2]).MappedDataSizes[1] = 5
	(*input[2]).MappedDataSizes[2] = 50

	output := getLocalityAwareReduceTaskSchedule(input)
	fmt.Println(output)

	if output[0] != 1 {
		panic("shuffle error")
	}
	if output[1] != 1 {
		panic("shuffle error")
	}
	if output[2] != 2 {
		panic("shuffle error")
	}
}
