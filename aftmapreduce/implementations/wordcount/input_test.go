package wordcount

import (
	"SDCC-Project/utility"
	"fmt"
	"testing"
)

const (
	testFilePath = "./test.txt"
)

func Test_MapService(t *testing.T) {

	rawInput := Input{10, "../../../test-input-data/input.txt"}

	output, err := rawInput.splitFile()
	if err != nil {
		panic(err)
	}

	for index := range output {
		fmt.Println("------------------------")
		fmt.Println(output[index])

	}
}

func TestReflection(t *testing.T) {

	data := Input{10, "test"}

	rawData := data.ToByte()

	deserializedData := Input{}

	if err := utility.Decode(rawData, &deserializedData); err != nil {
		panic(err)
	}

	fmt.Println(deserializedData.FileDigest)
	fmt.Println(deserializedData.MapCardinality)
}
