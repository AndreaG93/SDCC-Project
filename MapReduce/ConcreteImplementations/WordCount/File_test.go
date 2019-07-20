package WordCount

import (
	"fmt"
	"testing"
)

const (
	testFilePath = "./test.txt"
)

var err error

func Test_MapService(t *testing.T) {

	rawInput := File{10, "../../../test-input-data/input.txt"}

	output, err := rawInput.splitFile()
	if err != nil {
		panic(err)
	}

	for index := range output {
		fmt.Println("------------------------")
		fmt.Println(output[index])

	}
}
