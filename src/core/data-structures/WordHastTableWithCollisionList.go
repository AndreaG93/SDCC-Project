package data_structures

import (
	"container/list"
	"core/utility"
	"fmt"
	"strings"
)

type WordHashTableWithCollisionList struct {
	table               []*list.List
	tableSize           uint
	numberOfOccurrences uint
}

func BuildWordHashTableWithCollisionList(size uint) *WordHashTableWithCollisionList {

	output := new(WordHashTableWithCollisionList)

	output.table = make([]*list.List, size)
	output.tableSize = size
	output.numberOfOccurrences = 0

	for index := uint(0); index < size; index++ {
		output.table[index] = list.New()
	}

	return output
}

func (obj *WordHashTableWithCollisionList) Insert(word string) error {

	var index uint
	var err error
	var currentList *list.List

	if index, err = utility.GenerateArrayIndexFromString(word, obj.tableSize); err != nil {
		return err
	}

	currentList = obj.table[index]

	for e := currentList.Front(); e != nil; e = e.Next() {

		currentWordToken := e.Value.(*WordToken)

		if strings.Compare(currentWordToken.Word, word) == 0 {
			currentWordToken.Occurrences++
			return nil
		} else if strings.Compare(currentWordToken.Word, word) < 0 {

			currentList.InsertBefore(BuildWordToken(word), e)
			obj.numberOfOccurrences++
			return nil
		}
	}

	currentList.PushBack(BuildWordToken(word))
	obj.numberOfOccurrences++
	return nil
}

func (obj *WordHashTableWithCollisionList) Print() {

	var currentList *list.List

	for i := uint(0); i < obj.tableSize; i++ {

		currentList = obj.table[i]

		fmt.Printf(" --- Array position: %d --- \n", i)

		for e := currentList.Front(); e != nil; e = e.Next() {

			currentWordToken := e.Value.(*WordToken)
			fmt.Println(*currentWordToken)
		}
	}
}

func (obj *WordHashTableWithCollisionList) GetDigest() (string, error) {
	return utility.SHA512(obj.serialize())
}

func (obj *WordHashTableWithCollisionList) WriteOnLocalDisk() {

	var data []serializedWordHashTableWithCollisionList

	data = obj.serialize()

	writeSerializedWordHashTableWithCollisionListOnLocalDisk(data)
}

func (obj *WordHashTableWithCollisionList) getCollisionListAt(index uint) *list.List {
	return obj.table[index]
}

func (obj *WordHashTableWithCollisionList) serialize() []serializedWordHashTableWithCollisionList {

	var currentList *list.List
	var serializedOutput []serializedWordHashTableWithCollisionList
	var indexUsedForSerializedOutput uint

	serializedOutput = make([]serializedWordHashTableWithCollisionList, obj.numberOfOccurrences+1)
	indexUsedForSerializedOutput = 1

	serializedOutput[0].TableSizeOrIndex = obj.tableSize

	for i := uint(0); i < obj.tableSize; i++ {

		currentList = obj.table[i]

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
