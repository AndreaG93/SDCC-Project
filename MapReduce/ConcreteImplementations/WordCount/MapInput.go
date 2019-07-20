package WordCount

import (
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenHashTable"
	"SDCC-Project/utility"
	"strings"
)

type MapInput struct {
	Input          string
	MapCardinality uint
}

func (obj MapInput) PerformTask() (string, []byte, error) {

	var err error
	var data *WordTokenHashTable.WordTokenHashTable

	data = WordTokenHashTable.New(obj.MapCardinality)
	wordScanner := utility.BuildWordScannerFromString(obj.Input)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		if err = data.InsertWord(currentWord); err != nil {
			return "", nil, err
		}
	}

	rawData, err := data.Serialize()
	if err != nil {
		return "", nil, err
	}

	digest := utility.GenerateDigestUsingSHA512(rawData)

	return digest, rawData, nil
}
