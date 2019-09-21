package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"fmt"
)

type Reduce struct {
}

type ReduceInput struct {
	LocalDataDigest string
	ReduceWorkIndex int
}

type ReduceOutput struct {
	Digest string
	NodeId int

	MyInternetAddress string
	CPUUtilization    int
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	var err error

	process.GetLogger().PrintInfoLevelLabeledMessage(ReduceTaskName, fmt.Sprintf("Local data digest: %s -- Reduce work index: %d", input.LocalDataDigest, input.ReduceWorkIndex))

	digest, rawData := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex)

	utility.CheckError(process.GetDataRegistry().Set(digest, rawData))
	(*output).Digest = digest
	(*output).NodeId = process.GetPropertyAsInteger(property.NodeID)

	(*output).CPUUtilization, err = utility.GetCPUPercentageUtilizationAsInteger()
	utility.CheckError(err)

	(*output).MyInternetAddress = process.GetPropertyAsString(property.InternetAddress)

	return nil
}

func performReduceTask(localDataDigest string, reduceTaskIndex int) (string, []byte) {

	localWordTokenHashTable := WordTokenHashTable.Deserialize(process.GetDataRegistry().Get(localDataDigest))

	localWordTokenList := localWordTokenHashTable.GetWordTokenListAt(reduceTaskIndex)

	receivedDataDigest := GetGuidAssociation(localDataDigest)

	for _, digest := range receivedDataDigest {

		currentWordTokenList := WordTokenList.Deserialize(process.GetDataRegistry().Get(digest))
		localWordTokenList.Merge(currentWordTokenList)
	}

	digest, rawData := localWordTokenList.GetDigestAndSerializedData()

	return digest, rawData
}
