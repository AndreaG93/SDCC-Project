package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
)

type Reduce struct {
}

type ReduceInput struct {
	LocalDataDigest    string
	ReceivedDataDigest []string
	ReduceWorkIndex    int
}

type ReduceOutput struct {
	Digest string
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	digest, wordTokenList := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex, input.ReceivedDataDigest)

	node.GetCache().Set(digest, wordTokenList)
	(*output).Digest = digest

	return nil
}

func performReduceTask(localDataDigest string, reduceTaskIndex int, receivedDataDigest []string) (string, *WordTokenList.WordTokenList) {

	localWordTokenList := (node.GetCache().Get(localDataDigest)).(*WordTokenHashTable.WordTokenHashTable).GetWordTokenListAt(reduceTaskIndex)

	for _, digest := range receivedDataDigest {

		currentWordTokenList := (node.GetCache().Get(digest)).(*WordTokenList.WordTokenList)
		localWordTokenList.Merge(currentWordTokenList)
	}

	rawData, err := localWordTokenList.Serialize()
	utility.CheckError(err)
	digest := utility.GenerateDigestUsingSHA512(rawData)

	return digest, localWordTokenList
}
