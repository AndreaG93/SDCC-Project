package WordTokenListGroupSet

import (
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenHashTable"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenListGroup"
)

type WordTokenListGroupSet struct {
	set  []*WordTokenListGroup.WordTokenListGroup
	size uint
}

func New(input []*WordTokenHashTable.WordTokenHashTable) *WordTokenListGroupSet {

	output := new(WordTokenListGroupSet)

	(*output).size = uint(len(input))
	(*output).set = make([]*WordTokenListGroup.WordTokenListGroup, (*output).size)

	for index := uint(0); index < (*output).size; index++ {

		(*output).set[index] = WordTokenListGroup.NewFrom(input, index)
	}

	return output
}

func (obj *WordTokenListGroupSet) GetGroup(groupIndex uint) *WordTokenListGroup.WordTokenListGroup {

	return (*obj).set[groupIndex]
}