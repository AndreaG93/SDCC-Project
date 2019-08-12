package wordcount

import (
	"SDCC-Project/utility"
	"fmt"
	"testing"
)

func TestSerializationForReduceTaskOutput(t *testing.T) {

	input := new(AFTReduceTaskOutput)

	(*input).ReplayDigest = "test"
	(*input).NodeIdsWithCorrectResult = []int{1, 2, 3, 4, 5}
	(*input).IdGroup = 1

	rawData, _ := utility.Encode(input)

	output := AFTReduceTaskOutput{}

	utility.Decode(rawData, &output)

	fmt.Println(output)
}

func TestSerializationForMapTaskOutput(t *testing.T) {

	input := new(AFTMapTaskOutput)

	(*input).ReplayDigest = "test"
	(*input).NodeIdsWithCorrectResult = []int{1, 2, 3, 4, 5}
	(*input).IdGroup = 1
	(*input).MappedDataSizes = make(map[int]int)

	(*input).MappedDataSizes[4] = 1
	(*input).MappedDataSizes[5] = 1

	rawData, _ := utility.Encode(input)

	output := AFTMapTaskOutput{}

	utility.Decode(rawData, &output)

	fmt.Println(output)
}
