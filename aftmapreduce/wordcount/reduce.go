package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"net/rpc"
)

type Reduce struct {
}

type ReduceInput struct {
	LocalDataDigest string
	ReduceWorkIndex int
}

type ReduceOutput struct {
	Digest string
}

func (x *Reduce) Execute(input ReduceInput, output *ReduceOutput) error {

	digest, wordTokenList := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex)

	node.GetDataRegistry().Set(digest, wordTokenList)
	(*output).Digest = digest

	wordTokenList.Print()

	return nil
}

func performReduceTask(localDataDigest string, reduceTaskIndex int) (string, *WordTokenList.WordTokenList) {

	localWordTokenList := (node.GetDataRegistry().Get(localDataDigest)).(*WordTokenHashTable.WordTokenHashTable).GetWordTokenListAt(reduceTaskIndex)
	receivedDataDigest := node.GetDigestRegistry().GetAssociatedDigest(localDataDigest)

	for _, digest := range receivedDataDigest {

		currentWordTokenList := (node.GetDataRegistry().Get(digest)).(*WordTokenList.WordTokenList)
		localWordTokenList.Merge(currentWordTokenList)
	}

	rawData, err := localWordTokenList.Serialize()
	utility.CheckError(err)
	digest := utility.GenerateDigestUsingSHA512(rawData)

	return digest, localWordTokenList
}

func sendReduceTask(NodeIds []int, GroupId int, localDataDigest string, receiverReduceTaskId int) {

	internetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPCWithIdConstraints(GroupId, aftmapreduce.WordCountReduceTaskRPCBasePort, NodeIds)

	for _, sender := range internetAddresses {

		input := new(ReduceInput)
		output := new(ReduceOutput)

		(*input).ReduceWorkIndex = receiverReduceTaskId
		(*input).LocalDataDigest = localDataDigest

		worker, err := rpc.Dial("tcp", sender)
		utility.CheckError(err)

		err = worker.Call("Reduce.Execute", input, output)
		utility.CheckError(worker.Close())
		if err == nil {
			return
		}
	}
}
