package data_structures

import (
	"core/utility"
	"fmt"
)

type WordTokenHashTable struct {
	hashTable     []*WordTokenList
	hashTableSize uint
}

func BuildWordTokenHashTable(size uint) *WordTokenHashTable {

	output := new(WordTokenHashTable)

	(*output).hashTable = make([]*WordTokenList, size)
	(*output).hashTableSize = size

	for index := uint(0); index < size; index++ {
		(*output).hashTable[index] = BuildWordTokenList()
	}

	return output
}

func (obj *WordTokenHashTable) InsertWordToken(wordToken *WordToken) error {

	var index uint
	var err error
	var currentWordTokenList *WordTokenList

	if index, err = utility.GenerateArrayIndexFromString((*wordToken).Word, (*obj).hashTableSize); err != nil {
		return err
	}

	currentWordTokenList = (*obj).hashTable[index]
	(*currentWordTokenList).InsertWordToken(wordToken)

	return nil
}

func (obj *WordTokenHashTable) InsertWord(word string) error {

	return (*obj).InsertWordToken(BuildWordToken(word, 1))
}

func (obj *WordTokenHashTable) Print() {

	var currentList *WordTokenList

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentList = (*obj).hashTable[index]

		fmt.Printf(" --- Array position: %d --- \n", index)
		(*currentList).Print()
	}
}

func (obj *WordTokenHashTable) Serialize() WordTokenHashTableSerialized {

	var output []WordToken
	var currentWordTokenList *WordTokenList
	var currentWordToken *WordToken
	totalNumberOfWordToken := uint(0)

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]
		totalNumberOfWordToken += (*currentWordTokenList).length
	}

	output = make(WordTokenHashTableSerialized, totalNumberOfWordToken+1)
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

	return output
}

func (obj *WordTokenHashTable) extractWordTokenListAt() *WordTokenList {
	return nil
}
