package data_structures

import (
	"SDCC-Project-WorkerNode/src/core/utility"
	"encoding/gob"
	"os"
)

type WordTokenHashTableSerialized []WordToken

func ReadWordTokenHashTableSerializedFromLocalDisk(filePath string) (WordTokenHashTableSerialized, error) {

	var inputFile *os.File
	var err error

	if inputFile, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return nil, err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()

	decoder := gob.NewDecoder(inputFile)
	output := WordTokenHashTableSerialized{}

	err = decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	return output, nil
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
