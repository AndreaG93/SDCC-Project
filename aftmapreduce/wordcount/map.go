package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"strings"
)

type Map struct {
}

type MapInput struct {
	Text               string
	MappingCardinality int
}

type MapOutput struct {
	IdNode          int
	IdGroup         int
	ReplayDigest    string
	MappedDataSizes map[int]int
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	node.GetLogger().PrintInfoStartingTaskMessage(MapTaskName)

	digest, wordTokenHashTable, mappedDataSizes := performMapTask(input.Text, input.MappingCardinality)

	node.GetDataRegistry().Set(digest, wordTokenHashTable.Serialize())

	(*output).IdGroup = node.GetPropertyAsInteger(property.NodeGroupID)
	(*output).IdNode = node.GetPropertyAsInteger(property.NodeID)
	(*output).ReplayDigest = digest
	(*output).MappedDataSizes = mappedDataSizes

	if node.GetPropertyAsInteger(property.NodeID) == 1 || node.GetPropertyAsInteger(property.NodeID) == 4 {
		(*output).ReplayDigest = "dsadasdasdas"
	}

	node.GetLogger().PrintInfoCompleteTaskMessage(MapTaskName)

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

	rawData := outputData.Serialize()
	digest := utility.GenerateDigestUsingSHA512(rawData)

	for index := 0; index < mappingCardinality; index++ {
		mappedDataSizes[index] = outputData.GetWordTokenListAt(index).GetLength()
	}

	return digest, outputData, mappedDataSizes
}
