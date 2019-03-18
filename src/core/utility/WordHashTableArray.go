package utility

import (
	"bufio"

	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

type WordHashTableArray []map[string]uint

func BuildWordHashTableArray(size uint) WordHashTableArray {

	output := make(WordHashTableArray, size)

	for index := uint(0); index < size; index++ {
		output[index] = make(map[string]uint, 1000)
	}

	return output
}

func (obj WordHashTableArray) GenerateDigest() string {

	return ""
}

func (obj WordHashTableArray) InsertWordsFromScanner(wordScanner *bufio.Scanner) {

	for wordScanner.Scan() {

		currentWord := strings.ToLower(wordScanner.Text())

		index, errorGettingWordDigest := GenerateArrayIndexFromString(currentWord, uint(len(obj)))
		CheckError(errorGettingWordDigest)

		obj.InsertWordAt(currentWord, index)
	}
}

func (obj WordHashTableArray) InsertWordAt(word string, index uint) {

	nestedHashTable := obj[index]
	mWordOccurrences := nestedHashTable[word]

	if mWordOccurrences == 0 {
		nestedHashTable[word] = 1
	} else {
		nestedHashTable[word] = mWordOccurrences + 1
	}
}

type record struct {
	Index       uint
	Word        string
	Occurrences uint
}

func (obj WordHashTableArray) WriteOnLocalFile(filePath string) error {

	var outputFile *os.File
	var err error

	totalNumberOfRecords := uint(0)

	for i := uint(0); i < uint(len(obj)); i++ {
		totalNumberOfRecords += uint(len(obj[i]))
	}

	serializedStructure := make([]record, totalNumberOfRecords)

	for i, index := uint(0), uint(0); i < uint(len(obj)); i++ {
		for key, value := range obj[i] {
			serializedStructure[index] = record{i, key, value}
			index++
		}
	}

	if outputFile, err = os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer func() {
		CheckError(outputFile.Close())
	}()

	encoder := gob.NewEncoder(outputFile)

	if err = encoder.Encode(serializedStructure); err != nil {
		return err
	}

	return nil
}

func (obj WordHashTableArray) ReadFromLocalFile(filePath string) error {

	var inputFile *os.File
	var err error

	if inputFile, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return err
	}
	defer func() {
		CheckError(inputFile.Close())
	}()

	decoder := gob.NewDecoder(inputFile)
	p := []record{}

	err = decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)

	for i := uint(0); i < uint(len(p)); i++ {

		var record record
		var nestedMap map[string]uint

		record = p[i]
		nestedMap = obj[record.Index]
		nestedMap[record.Word] = record.Occurrences
	}

	return nil

}

func (obj WordHashTableArray) Print() {

	for i := uint(0); i < uint(len(obj)); i++ {

		fmt.Printf("Array position: %d\n", i)
		fmt.Println(obj[i])
	}
}
