package DataStructures

import (
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenHashTable"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenList"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenListGroup"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenListGroupSet"
	"SDCC-Project/utility"
	"testing"
)

const (
	totalHashTables = 5
)

func getPopulatedHashTable() *WordTokenHashTable.WordTokenHashTable {

	output := WordTokenHashTable.New(10)

	if err := (*output).InsertWord("Andrea"); err != nil {
		panic(err)
	}
	if err := (*output).InsertWord("Akko"); err != nil {
		panic(err)
	}
	if err := (*output).InsertWord("Yumi"); err != nil {
		panic(err)
	}
	if err := (*output).InsertWord("Hanabi"); err != nil {
		panic(err)
	}

	return output
}

func Test_Serialize(t *testing.T) {

	finalOutput := WordTokenList.New()
	output := make([]*WordTokenList.WordTokenList, totalHashTables)
	hashTables := make([]*WordTokenHashTable.WordTokenHashTable, totalHashTables)

	for index := 0; index < 5; index++ {
		hashTables[index] = getPopulatedHashTable()
	}

	set := WordTokenListGroupSet.New(hashTables)

	for index := 0; index < 5; index++ {

		group := set.GetGroup(uint(index))

		rawGroup, err := group.Serialize()
		utility.CheckError(err)

		reserializedGroup, err := WordTokenListGroup.Deserialize(rawGroup)
		utility.CheckError(err)

		output[index] = reserializedGroup.Merge()
	}

	for _, list := range output {
		finalOutput.Merge(list)
	}

	finalOutput.Print()
}
