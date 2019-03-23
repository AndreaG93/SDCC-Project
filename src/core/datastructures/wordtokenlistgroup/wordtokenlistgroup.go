package wordtokenlistgroup

import (
	"SDCC-Project-WorkerNode/src/core/datastructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/src/core/datastructures/wordtokenlist"
	"SDCC-Project-WorkerNode/src/core/utility"
	"fmt"
)

type WordTokenListGroup struct {
	data      []*wordtokenlist.WordTokenList
	groupSize uint
}

type wordTokenListGroupUnit struct {
	indexGroup  uint
	Word        string
	Occurrences uint
}

func New(input []*wordtokenhashtable.WordTokenHashTable, groupIndex uint) *WordTokenListGroup {

	size := uint(len(input))

	output := new(WordTokenListGroup)
	(*output).groupSize = 0
	(*output).data = make([]*wordtokenlist.WordTokenList, size)

	for index := uint(0); index < size; index++ {
		(*output).data[index] = input[index].GetWordTokenListAt(groupIndex)
		(*output).groupSize += (*output).data[index].GetLength()
	}

	return output
}

func (obj *WordTokenListGroup) Serialize() ([]byte, error) {

	return utility.Encode(output)
}

func (obj *WordTokenListGroup) Print() {

	var currentList *wordtokenlist.WordTokenList

	for index := uint(0); index < uint(len((*obj).data)); index++ {

		currentList = (*obj).data[index]

		fmt.Printf(" --- Position: %d --- \n", index)
		(*currentList).Print()
	}
}
