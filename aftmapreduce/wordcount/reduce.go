package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
)

type Reduce struct {
}

type ReduceInput struct {
	LocalDataDigest string
	ReduceWorkIndex int
}

type ReduceOutput struct {
	Digest string
	nodeId int
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	digest, rawData := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex)

	node.GetDataRegistry().Set(digest, rawData)
	(*output).Digest = digest
	(*output).nodeId = node.GetPropertyAsInteger(property.NodeID)

	return nil
}

func performReduceTask(localDataDigest string, reduceTaskIndex int) (string, []byte) {

	localWordTokenList := (node.GetDataRegistry().Get(localDataDigest)).(*WordTokenHashTable.WordTokenHashTable).GetWordTokenListAt(reduceTaskIndex)
	receivedDataDigest := node.GetDigestRegistry().GetAssociatedDigest(localDataDigest)

	for _, digest := range receivedDataDigest {

		currentWordTokenList := (node.GetDataRegistry().Get(digest)).(*WordTokenList.WordTokenList)
		localWordTokenList.Merge(currentWordTokenList)
	}

	digest, rawData := localWordTokenList.GetDigestAndSerializedData()

	return digest, rawData
}
