package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/utility"
	"strings"
)

type Map struct {
}

type MapInput struct {
	text               string
	MappingCardinality int
}

type MapOutput struct {
	IdNode          int
	IdGroup         int
	ReplayDigest    string
	MappedDataSizes map[int]int
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	digest, wordTokenHashTable, mappedDataSizes := performMapTask(input.text, input.MappingCardinality)

	node.GetCache().Set(digest, wordTokenHashTable)

	(*output).IdGroup = node.GetIntegerProperty(property.NodeID)
	(*output).IdNode = node.GetIntegerProperty(property.NodeGroupID)
	(*output).ReplayDigest = digest
	(*output).MappedDataSizes = mappedDataSizes

	return nil
}

func performMapTask(text string, mappingCardinality int) (string, *WordTokenHashTable.WordTokenHashTable, map[int]int) {

	mappedDataSizes := make(map[int]int)

	outputData := WordTokenHashTable.New(uint(mappingCardinality))
	wordScanner := utility.BuildWordScannerFromString(text)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		utility.CheckError(outputData.InsertWord(currentWord))
	}

	rawData, err := outputData.Serialize()
	utility.CheckError(err)
	digest := utility.GenerateDigestUsingSHA512(rawData)

	for index := 0; index < mappingCardinality; index++ {
		mappedDataSizes[index] = outputData.GetWordTokenListAt(index).GetLength()
	}

	return digest, outputData, mappedDataSizes
}
