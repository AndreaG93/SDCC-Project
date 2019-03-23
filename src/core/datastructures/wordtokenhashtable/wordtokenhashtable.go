package wordtokenhashtable

import (
	"SDCC-Project-WorkerNode/src/core/datastructures/wordtoken"
	"SDCC-Project-WorkerNode/src/core/datastructures/wordtokenlist"
	"SDCC-Project-WorkerNode/src/core/utility"
	"fmt"
)

type WordTokenHashTable struct {
	hashTable     []*wordtokenlist.WordTokenList
	hashTableSize uint
}

func New(size uint) *WordTokenHashTable {

	output := new(WordTokenHashTable)

	(*output).hashTable = make([]*wordtokenlist.WordTokenList, size)
	(*output).hashTableSize = size

	for index := uint(0); index < size; index++ {
		(*output).hashTable[index] = wordtokenlist.New()
	}

	return output
}

func Deserialize(input []byte) (*WordTokenHashTable, error) {

	var output *WordTokenHashTable
	var currentWordToken *wordtoken.WordToken

	serializedData := []wordtoken.WordToken{}

	if err := utility.Decode(input, &serializedData); err != nil {
		return nil, err
	}

	output = New(serializedData[0].Occurrences)

	for index := uint(1); index < uint(len(serializedData)); index++ {

		currentWordToken = wordtoken.New(serializedData[index].Word, serializedData[index].Occurrences)
		if err := (*output).InsertWordToken(currentWordToken); err != nil {
			return nil, err
		}

	}

	return output, nil
}

func (obj *WordTokenHashTable) InsertWordToken(wordToken *wordtoken.WordToken) error {

	var index uint
	var err error
	var currentWordTokenList *wordtokenlist.WordTokenList

	if index, err = utility.GenerateArrayIndexFromString((*wordToken).Word, (*obj).hashTableSize); err != nil {
		return err
	}

	currentWordTokenList = (*obj).hashTable[index]
	(*currentWordTokenList).InsertWordToken(wordToken)

	return nil
}

func (obj *WordTokenHashTable) InsertWord(word string) error {

	return (*obj).InsertWordToken(wordtoken.New(word, 1))
}

func (obj *WordTokenHashTable) Print() {

	var currentList *wordtokenlist.WordTokenList

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentList = (*obj).hashTable[index]

		fmt.Printf(" --- Array position: %d --- \n", index)
		(*currentList).Print()
	}
}

func (obj *WordTokenHashTable) GetWordTokenListAt(index uint) *wordtokenlist.WordTokenList {
	return obj.hashTable[index]
}

func (obj *WordTokenHashTable) Serialize() ([]byte, error) {

	var output []wordtoken.WordToken
	var currentWordTokenList *wordtokenlist.WordTokenList
	var currentWordToken *wordtoken.WordToken
	totalNumberOfWordToken := uint(0)

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]
		totalNumberOfWordToken += (*currentWordTokenList).GetLength()
	}

	output = make([]wordtoken.WordToken, totalNumberOfWordToken+1)
	output[0].Word = ""
	output[0].Occurrences = (*obj).hashTableSize

	for index, outputIndex := uint(0), uint(1); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]

		(*currentWordTokenList).IteratorReset()

		for (*currentWordTokenList).Next() {

			currentWordToken = (*currentWordTokenList).WordToken()

			output[outputIndex].Word = (*currentWordToken).Word
			output[outputIndex].Occurrences = (*currentWordToken).Occurrences

			outputIndex++
		}
	}

	return utility.Encode(output)
}
