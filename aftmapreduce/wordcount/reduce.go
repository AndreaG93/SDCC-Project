package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
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
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	node.GetLogger().PrintInfoTaskMessage(ReduceTaskName, fmt.Sprintf("Local data digest: %s -- Reduce work index: %d", input.LocalDataDigest, input.ReduceWorkIndex))

	digest, rawData := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex)

	utility.CheckError(node.GetDataRegistry().Set(digest, rawData))
	(*output).Digest = digest
	(*output).NodeId = node.GetPropertyAsInteger(property.NodeID)

	return nil
}

func performReduceTask(localDataDigest string, reduceTaskIndex int) (string, []byte) {

	localWordTokenHashTable := WordTokenHashTable.Deserialize(node.GetDataRegistry().Get(localDataDigest))

	localWordTokenList := localWordTokenHashTable.GetWordTokenListAt(reduceTaskIndex)

	receivedDataDigest := GetGuidAssociation(localDataDigest)

	for _, digest := range receivedDataDigest {

		currentWordTokenList := WordTokenList.Deserialize(node.GetDataRegistry().Get(digest))
		localWordTokenList.Merge(currentWordTokenList)
	}

	digest, rawData := localWordTokenList.GetDigestAndSerializedData()

	return digest, rawData
}
