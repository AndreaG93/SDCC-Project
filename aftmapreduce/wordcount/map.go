package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
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

	MyInternetAddress string
	CPUUtilization    int
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	var err error
	var digest string
	var indexPartitionSize map[int]int
	var wordTokenHashTable *WordTokenHashTable.WordTokenHashTable
	var rawData []byte

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Received %s task! -- Mapping Cardinality %d", MapTaskName, input.MappingCardinality))

	guid := utility.GenerateDigestUsingSHA512([]byte(input.Text))

	if rawData = process.GetDataRegistry().Get(guid); rawData == nil {

		digest, wordTokenHashTable, indexPartitionSize = performMapTask(input.Text, input.MappingCardinality)

		if err = process.GetDataRegistry().Set(digest, wordTokenHashTable.Serialize()); err != nil {
			return err
		}
		if err = process.GetDataRegistry().Set(guid, []byte(digest)); err != nil {
			return err
		}
		if err = process.GetDataRegistry().Set(fmt.Sprintf("%s-mappedDataSize", guid), utility.Encode(indexPartitionSize)); err != nil {
			return err
		}

	} else {

		digest = string(rawData)
		rawData = process.GetDataRegistry().Get(fmt.Sprintf("%s-mappedDataSize", guid))

		if err = utility.Decoding(rawData, &indexPartitionSize); err != nil {
			return err
		}
	}

	if (*output).CPUUtilization, err = utility.GetCPUPercentageUtilizationAsInteger(); err != nil {
		return err
	}

	(*output).MyInternetAddress = process.GetPropertyAsString(property.InternetAddress)
	(*output).IdGroup = process.GetPropertyAsInteger(property.NodeGroupID)
	(*output).IdNode = process.GetPropertyAsInteger(property.NodeID)
	(*output).ReplayDigest = digest
	(*output).MappedDataSizes = indexPartitionSize

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("A Map task is COMPLETED with Digest %s :: IndexPartitionSize %d", (*output).ReplayDigest, (*output).MappedDataSizes))

	return nil
}

func performMapTask(text string, mappingCardinality int) (string, *WordTokenHashTable.WordTokenHashTable, map[int]int) {

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Test received: %s", text))

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
