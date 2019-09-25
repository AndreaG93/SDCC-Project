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
	var digestData string
	var indexPartitionSizes map[int]int
	var serializedData []byte

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Received %s task! -- Mapping Cardinality %d", MapTaskName, input.MappingCardinality))

	guid := utility.GenerateDigestUsingSHA512([]byte(input.Text))

	if serializedData = process.GetDataRegistry().Get(guid); serializedData == nil {

		if digestData, serializedData, indexPartitionSizes, err = performMapTask(input.Text, input.MappingCardinality); err != nil {
			return err
		} else {

			if serializedIndexPartitionSize, err := utility.Encoding(indexPartitionSizes); err != nil {
				return err
			} else {

				if err = process.GetDataRegistry().Set(digestData, serializedData); err != nil {
					return err
				}
				if err = process.GetDataRegistry().Set(fmt.Sprintf("%s-mappedDataSize", guid), serializedIndexPartitionSize); err != nil {
					return err
				}
			}
		}

	} else {

		digestData = string(serializedData)
		serializedData = process.GetDataRegistry().Get(fmt.Sprintf("%s-mappedDataSize", guid))

		if err = utility.Decoding(serializedData, &indexPartitionSizes); err != nil {
			return err
		}
	}

	if (*output).CPUUtilization, err = utility.GetCPUPercentageUtilizationAsInteger(); err != nil {
		return err
	}

	(*output).MyInternetAddress = process.GetPropertyAsString(property.InternetAddress)
	(*output).IdGroup = process.GetPropertyAsInteger(property.NodeGroupID)
	(*output).IdNode = process.GetPropertyAsInteger(property.NodeID)
	(*output).ReplayDigest = digestData
	(*output).MappedDataSizes = indexPartitionSizes

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("A Map task is COMPLETED with Digest %s :: IndexPartitionSize %d", (*output).ReplayDigest, (*output).MappedDataSizes))

	crash()

	if isOccurredAnArbitraryCrash() {
		(*output).ReplayDigest = "dsadsadasfewhfiuehuiyw7i34yt78f2g3g7823gafsf43fwet34ghdhdgrsgt--GGGGGGGGGG"
	}

	return nil
}

func performMapTask(text string, mappingCardinality int) (string, []byte, map[int]int, error) {

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Text received: %s", text))

	indexPartitionSizes := make(map[int]int)

	outputData := WordTokenHashTable.New(uint(mappingCardinality))
	wordScanner := utility.BuildWordScannerFromString(text)

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())
		utility.CheckError(outputData.InsertWord(currentWord))
	}

	if digestData, serializedData, err := outputData.GetDigestAndSerializedData(); err != nil {
		return "", nil, nil, err
	} else {

		for index := 0; index < mappingCardinality; index++ {
			indexPartitionSizes[index] = outputData.GetWordTokenListAt(index).GetLength()
		}

		return digestData, serializedData, indexPartitionSizes, nil
	}
}
