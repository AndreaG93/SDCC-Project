package wordtokenlistgroupset

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenlistgroup"
)

type WordTokenListGroupSet struct {
	set  []*wordtokenlistgroup.WordTokenListGroup
	size uint
}

func New(input []*wordtokenhashtable.WordTokenHashTable) *WordTokenListGroupSet {

	output := new(WordTokenListGroupSet)

	(*output).size = uint(len(input))
	(*output).set = make([]*wordtokenlistgroup.WordTokenListGroup, (*output).size)

	for index := uint(0); index < (*output).size; index++ {

		(*output).set[index] = wordtokenlistgroup.NewFrom(input, index)

	}

	return output
}

func (obj *WordTokenListGroupSet) GetGroup(groupIndex uint) *wordtokenlistgroup.WordTokenListGroup {

	return (*obj).set[groupIndex]
}
