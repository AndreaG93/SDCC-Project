package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenhashtable"
	"io/ioutil"
	"testing"
)

const (
	testFilePath = "./test.txt"
)

var err error

func Test_MapService(t *testing.T) {

	var rawData []byte

	if rawData, err = ioutil.ReadFile(testFilePath); err != nil {
		panic(err)
	}

	digest := mapTest(string(rawData))
	mapGetTest(digest)
}

func mapTest(text string) string {

	input := MapInput{text, 5}
	output := MapOutput{}
	object := Map{}

	if err = object.Execute(input, &output); err != nil {
		panic(err)
	}

	return output.digest
}

func mapGetTest(digest string) {

	var hashTable *wordtokenhashtable.WordTokenHashTable

	input := MapGetInput{digest}
	output := MapGetOutput{}

	object := MapGet{}

	if err = object.Execute(input, &output); err != nil {
		panic(err)
	}

	if hashTable, err = wordtokenhashtable.Deserialize(output.data); err != nil {
		panic(err)
	}

	hashTable.Print()
}
