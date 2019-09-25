package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
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

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Start a Reduce Task involving a local data with digest: %s -- Reduce work index: %d", input.LocalDataDigest, input.ReduceWorkIndex))

	if digestData, serializedData, err := performReduceTask(input.LocalDataDigest, input.ReduceWorkIndex); err != nil {
		return err
	} else {

		if err = process.GetDataRegistry().Set(digestData, serializedData); err != nil {
			return err
		} else {
			(*output).Digest = digestData
			(*output).NodeId = process.GetPropertyAsInteger(property.NodeID)

			return nil
		}
	}
}

func performReduceTask(localDataDigest string, reduceTaskIndex int) (string, []byte, error) {

	if localWordTokenHashTable, err := WordTokenHashTable.Deserialize(process.GetDataRegistry().Get(localDataDigest)); err != nil {
		return "", nil, err
	} else {
		localWordTokenList := localWordTokenHashTable.GetWordTokenListAt(reduceTaskIndex)

		if digestsOfReceivedPartitions, err := GetDigestAssociationArray(localDataDigest, reduceTaskIndex); err != nil {
			return "", nil, err
		} else {

			for _, digest := range digestsOfReceivedPartitions {

				if currentWordTokenList, err := WordTokenList.Deserialize(process.GetDataRegistry().Get(digest)); err != nil {
					return "", nil, err
				} else {
					localWordTokenList.Merge(currentWordTokenList)
				}
			}

			return localWordTokenList.GetDigestAndSerializedData()
		}
	}
}
