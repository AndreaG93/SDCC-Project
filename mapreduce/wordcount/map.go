package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/DataStructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/Register/Worker/WorkerMapRegister"
	"SDCC-Project-WorkerNode/utility"
	"strings"
)

type Map struct {
}

type MapInput struct {
	Input          string
	MapCardinality uint
}

type MapOutput struct {
	Digest string
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	var err error
	var data *wordtokenhashtable.WordTokenHashTable

	data = wordtokenhashtable.New(input.MapCardinality)

	wordScanner := utility.BuildWordScannerFromString(input.Input)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = data.InsertWord(currentWord); err != nil {
			return err
		}
	}

	rawData, err := data.Serialize()
	if err != nil {
		return err
	}

	digest := utility.GenerateDigestUsingSHA512(rawData)

	WorkerMapRegister.GetInstance().Set(digest, rawData)

	output.Digest = digest

	return nil
}
