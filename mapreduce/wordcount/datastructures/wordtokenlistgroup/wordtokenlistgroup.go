package wordtokenlistgroup

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtoken"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenhashtable"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtokenlist"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
)

type WordTokenListGroup struct {
	group           []*wordtokenlist.WordTokenList
	size            uint
	totalPopulation uint
}

type wordTokenListGroupSerializationUnit struct {
	word                            string
	occurrences                     uint
	membershipListIndexAndGroupSize uint
}

func New(size uint) *WordTokenListGroup {

	output := new(WordTokenListGroup)

	(*output).totalPopulation = 0
	(*output).size = uint(size)
	(*output).group = make([]*wordtokenlist.WordTokenList, (*output).size)

	for index := uint(0); index < (*output).size; index++ {
		(*output).group[0] = wordtokenlist.New()
	}

	return output
}

func NewFrom(input []*wordtokenhashtable.WordTokenHashTable, groupIndex uint) *WordTokenListGroup {

	output := new(WordTokenListGroup)

	(*output).totalPopulation = 0
	(*output).size = uint(len(input))
	(*output).group = make([]*wordtokenlist.WordTokenList, (*output).size)

	for index := uint(0); index < (*output).size; index++ {
		(*output).group[index] = input[index].GetWordTokenListAt(groupIndex)
		(*output).totalPopulation += (*output).group[index].GetLength()
	}

	return output
}

func Deserialize(input []byte) (*WordTokenListGroup, error) {

	var output *WordTokenListGroup
	var currentWordToken *wordtoken.WordToken
	var currentWordTokenList *wordtokenlist.WordTokenList

	serializedData := []wordTokenListGroupSerializationUnit{}

	if err := utility.Decode(input, &serializedData); err != nil {
		return nil, err
	}

	output = New(serializedData[0].membershipListIndexAndGroupSize)

	for index := uint(1); index < uint(len(serializedData)); index++ {

		currentWordTokenList = output.group[serializedData[index].membershipListIndexAndGroupSize]
		currentWordToken = wordtoken.New(serializedData[index].word, serializedData[index].occurrences)

		currentWordTokenList.InsertWordToken(currentWordToken)
	}

	return output, nil
}

func (obj *WordTokenListGroup) Serialize() ([]byte, error) {

	var currentWordToken *wordtoken.WordToken
	var currentWordTokenList *wordtokenlist.WordTokenList
	var output []wordTokenListGroupSerializationUnit
	var outputIndex uint

	outputIndex = uint(0)
	output = make([]wordTokenListGroupSerializationUnit, (*obj).totalPopulation+1)

	output[0].membershipListIndexAndGroupSize = (*obj).size

	for index := uint(1); index < (*obj).size; index++ {

		currentWordTokenList = (*obj).group[index]

		for currentWordTokenList.Next() {

			currentWordToken = currentWordTokenList.WordToken()

			output[outputIndex].occurrences = (*currentWordToken).Occurrences
			output[outputIndex].word = (*currentWordToken).Word
			output[outputIndex].membershipListIndexAndGroupSize = index
		}
	}

	return utility.Encode(output)
}

func (obj *WordTokenListGroup) Merge() *wordtokenlist.WordTokenList {

	output := wordtokenlist.New()

	for index := uint(0); index < (*obj).size; index++ {

		currentWordTokenList := (*obj).group[index]
		(*output).Merge(currentWordTokenList)
	}

	return output
}

func (obj *WordTokenListGroup) Print() {

	var currentList *wordtokenlist.WordTokenList

	for index := uint(0); index < uint(len((*obj).group)); index++ {

		currentList = (*obj).group[index]

		fmt.Printf(" --- Position: %d --- \n", index)
		(*currentList).Print()
	}
}
