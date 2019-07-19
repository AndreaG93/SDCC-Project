package WordTokenListGroup

import (
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordToken"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenHashTable"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"fmt"
)

type WordTokenListGroup struct {
	group           []*WordTokenList.WordTokenList
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
	(*output).group = make([]*WordTokenList.WordTokenList, (*output).size)

	for index := uint(0); index < (*output).size; index++ {
		(*output).group[0] = WordTokenList.New()
	}

	return output
}

func NewFrom(input []*WordTokenHashTable.WordTokenHashTable, groupIndex uint) *WordTokenListGroup {

	output := new(WordTokenListGroup)

	(*output).totalPopulation = 0
	(*output).size = uint(len(input))
	(*output).group = make([]*WordTokenList.WordTokenList, (*output).size)

	for index := uint(0); index < (*output).size; index++ {
		(*output).group[index] = input[index].GetWordTokenListAt(groupIndex)
		(*output).totalPopulation += (*output).group[index].GetLength()
	}

	return output
}

func Deserialize(input []byte) (*WordTokenListGroup, error) {

	var output *WordTokenListGroup
	var currentWordToken *WordToken.WordToken
	var currentWordTokenList *WordTokenList.WordTokenList

	serializedData := []wordTokenListGroupSerializationUnit{}

	if err := utility.Decode(input, &serializedData); err != nil {
		return nil, err
	}

	output = New(serializedData[0].membershipListIndexAndGroupSize)

	for index := uint(1); index < uint(len(serializedData)); index++ {

		currentWordTokenList = output.group[serializedData[index].membershipListIndexAndGroupSize]
		currentWordToken = WordToken.New(serializedData[index].word, serializedData[index].occurrences)

		currentWordTokenList.InsertWordToken(currentWordToken)
	}

	return output, nil
}

func (obj *WordTokenListGroup) Serialize() ([]byte, error) {

	var currentWordToken *WordToken.WordToken
	var currentWordTokenList *WordTokenList.WordTokenList
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

func (obj *WordTokenListGroup) Merge() *WordTokenList.WordTokenList {

	output := WordTokenList.New()

	for index := uint(0); index < (*obj).size; index++ {

		currentWordTokenList := (*obj).group[index]
		(*output).Merge(currentWordTokenList)
	}

	return output
}

func (obj *WordTokenListGroup) Print() {

	var currentList *WordTokenList.WordTokenList

	for index := uint(0); index < uint(len((*obj).group)); index++ {

		currentList = (*obj).group[index]

		fmt.Printf(" --- Position: %d --- \n", index)
		(*currentList).Print()
	}
}
