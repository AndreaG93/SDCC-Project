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

	inputDigest := utility.GenerateDigestUsingSHA512([]byte(input.Text))

	if node.GetDataRegistry().Get(inputDigest) == nil {

		digest, wordTokenHashTable, mappedDataSizes := performMapTask(input.Text, input.MappingCardinality)

		node.GetDataRegistry().Set(digest, wordTokenHashTable.Serialize())

		(*output).ReplayDigest = digest
		(*output).MappedDataSizes = mappedDataSizes

		node.GetDataRegistry().Set(inputDigest, []byte((*output).ReplayDigest))
		node.GetDataRegistry().Set(inputDigest+"map", utility.Encode(mappedDataSizes))

		node.GetLogger().PrintInfoCompleteTaskMessage(MapTaskName)

	} else {

		(*output).ReplayDigest = string(node.GetDataRegistry().Get(inputDigest))

		mappedDataSizes := map[int]int{}
		utility.Decode(node.GetDataRegistry().Get(inputDigest+"map"), &mappedDataSizes)

		(*output).MappedDataSizes = mappedDataSizes

	}

	(*output).IdGroup = node.GetPropertyAsInteger(property.NodeGroupID)
	(*output).IdNode = node.GetPropertyAsInteger(property.NodeID)

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
