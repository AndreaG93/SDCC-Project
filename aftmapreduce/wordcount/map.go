package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"fmt"
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

	guid := utility.GenerateDigestUsingSHA512([]byte(input.Text))

	if !isMapTaskRequestDuplicated(guid) {

		digest, wordTokenHashTable, mappedDataSizes := performMapTask(input.Text, input.MappingCardinality)

		node.GetDataRegistry().Set(digest, wordTokenHashTable.Serialize())
		node.GetDataRegistry().Set(guid, []byte(digest))
		node.GetDataRegistry().Set(fmt.Sprintf("%s-mappedDataSize", guid), utility.Encode(mappedDataSizes))

		(*output).ReplayDigest = digest
		(*output).MappedDataSizes = mappedDataSizes

		node.GetLogger().PrintInfoCompleteTaskMessage(MapTaskName)

	} else {

		(*output).ReplayDigest = string(node.GetDataRegistry().Get(guid))
		utility.Decode(node.GetDataRegistry().Get(fmt.Sprintf("%s-mappedDataSize", guid)), &output.MappedDataSizes)
	}

	(*output).IdGroup = node.GetPropertyAsInteger(property.NodeGroupID)
	(*output).IdNode = node.GetPropertyAsInteger(property.NodeID)

	return nil
}

func isMapTaskRequestDuplicated(guid string) bool {
	return node.GetDataRegistry().Get(guid) != nil
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
