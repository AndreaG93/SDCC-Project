package WordTokenHashTable

import (
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordToken"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"fmt"
)

type WordTokenHashTable struct {
	hashTable     []*WordTokenList.WordTokenList
	hashTableSize uint
}

func New(size uint) *WordTokenHashTable {

	output := new(WordTokenHashTable)

	(*output).hashTable = make([]*WordTokenList.WordTokenList, size)
	(*output).hashTableSize = size

	for index := uint(0); index < size; index++ {
		(*output).hashTable[index] = WordTokenList.New()
	}

	return output
}

func (obj *WordTokenHashTable) InsertWordToken(wordToken *WordToken.WordToken) error {

	var index uint
	var err error
	var currentWordTokenList *WordTokenList.WordTokenList

	if index, err = utility.GenerateArrayIndexFromString((*wordToken).Word, (*obj).hashTableSize); err != nil {
		return err
	}

	currentWordTokenList = (*obj).hashTable[index]
	(*currentWordTokenList).InsertWordToken(wordToken)

	return nil
}

func (obj *WordTokenHashTable) InsertWord(word string) error {

	return (*obj).InsertWordToken(WordToken.New(word, 1))
}

func (obj *WordTokenHashTable) Print() {

	var currentList *WordTokenList.WordTokenList

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentList = (*obj).hashTable[index]

		fmt.Printf(" --- Array position: %d --- \n", index)
		(*currentList).Print()
	}
}

func (obj *WordTokenHashTable) GetWordTokenListAt(index int) *WordTokenList.WordTokenList {
	return obj.hashTable[index]
}

func (obj *WordTokenHashTable) Serialize() ([]byte, error) {

	var output []WordToken.WordToken
	var currentWordTokenList *WordTokenList.WordTokenList
	var currentWordToken *WordToken.WordToken
	totalNumberOfWordToken := 0

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]
		totalNumberOfWordToken += (*currentWordTokenList).GetLength()
	}

	output = make([]WordToken.WordToken, totalNumberOfWordToken+1)
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
