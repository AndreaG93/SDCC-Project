package data_structures

import (
	"container/list"
	"core/utility"
	"fmt"
)

type WordTokenHashTableSerialized []WordToken

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
	totalNumberOfWordToken := uint(0)

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]
		totalNumberOfWordToken += (*currentWordTokenList).length
	}

	output = make(WordTokenHashTableSerialized, totalNumberOfWordToken+1)
	output[0].Word = ""
	output[0].Occurrences = (*obj).hashTableSize

	for index := uint(0); index < (*obj).hashTableSize; index++ {

		currentWordTokenList = (*obj).hashTable[index]
		(*currentWordTokenList).IteratorReset()

		for (*currentWordTokenList).Next() {

			output[index+1].Word = ((*currentWordTokenList).WordToken()).Word
			output[index+1].Occurrences = ((*currentWordTokenList).WordToken()).Occurrences
		}
	}

	return output
}

func (obj WordTokenHashTableSerialized) Deserialize() (*WordTokenHashTable, error) {

	var output *WordTokenHashTable
	var currentWordToken *WordToken

	output = BuildWordTokenHashTable(obj[0].Occurrences)

	for index := uint(1); index < uint(len(obj)); index++ {

		currentWordToken = BuildWordToken(obj[index].Word, obj[index].Occurrences)
		if err := (*output).InsertWordToken(currentWordToken); err != nil {
			return nil, err
		}

	}

	return output, nil
}

func (obj *WordTokenHashTable) GetDigest() (string, error) {
	return utility.SHA512(obj.serialize())
}

func (obj *WordTokenHashTable) WriteOnLocalDisk() error {

	var serializedData TokenWordHashTableSerialized
	var err error

	fmt.Println(utility.SHA512(obj.serialize()))

	serializedData = obj.serialize()

	if err = serializedData.writeOnLocalDisk(); err != nil {
		return err
	}

	return nil
}

func ReadTokenWordHashTableFromLocalDisk(pathFile string) (*WordTokenHashTable, error) {

	var serializedData []TokenWordHashTableSerializationUnit
	var err error

	if serializedData, err = readSerializedWordHashTableWithCollisionListFromLocalDisk(pathFile); err != nil {
		return nil, err
	}

	return Deserialize(serializedData), nil
}

func (obj *WordTokenHashTable) getCollisionListAt(index uint) *list.List {
	return obj.hashTable[index]
}

func (obj *WordTokenHashTable) serialize() TokenWordHashTableSerialized {

	var currentList *list.List
	var serializedOutput []TokenWordHashTableSerializationUnit
	var indexUsedForSerializedOutput uint

	serializedOutput = make(TokenWordHashTableSerialized, obj.numberOfDistinctWords+1)
	indexUsedForSerializedOutput = 1

	serializedOutput[0].TableSizeOrIndex = obj.hashTableSize

	for i := uint(0); i < obj.hashTableSize; i++ {

		currentList = obj.hashTable[i]

		for e := currentList.Front(); e != nil; e = e.Next() {

			currentWordToken := e.Value.(*WordToken)

			serializedOutput[indexUsedForSerializedOutput].TableSizeOrIndex = i
			serializedOutput[indexUsedForSerializedOutput].Word = currentWordToken.Word
			serializedOutput[indexUsedForSerializedOutput].Occurrences = currentWordToken.Occurrences

			indexUsedForSerializedOutput++
		}
	}

	return serializedOutput
}

func Deserialize(input []TokenWordHashTableSerializationUnit) *WordTokenHashTable {

	var output *WordTokenHashTable
	var currentWordToken *WordToken

	outputHashTableSize := input[0].TableSizeOrIndex
	output = BuildWordTokenHashTable(outputHashTableSize)

	for index := uint(1); index < uint(len(input)); index++ {

		currentWord := input[index].Word
		currentWordOccurrences := input[index].Occurrences

		currentWordToken = BuildWordToken(currentWord, currentWordOccurrences)

		output.InsertWordToken(currentWordToken)
	}

	return output
}
