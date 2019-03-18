package data_structures

import (
	"core/utility"
	"encoding/gob"
	"fmt"
	"os"
)

type serializedWordHashTableWithCollisionList struct {
	TableSizeOrIndex uint
	Word             string
	Occurrences      uint
}

func writeSerializedWordHashTableWithCollisionListOnLocalDisk(data []serializedWordHashTableWithCollisionList) error {

	var outputFile *os.File
	var outputFileName string
	var err error

	if outputFileName, err = utility.SHA512(data); err != nil {
		return err
	}
	if outputFile, err = os.OpenFile(outputFileName, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(outputFile.Close())
	}()

	encoder := gob.NewEncoder(outputFile)

	if err = encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func ReadSerializedWordHashTableWithCollisionListFromLocalDisk(filePath string) ([]serializedWordHashTableWithCollisionList, error) {

	var inputFile *os.File
	var err error

	if inputFile, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return nil, err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()

	decoder := gob.NewDecoder(inputFile)
	output := []serializedWordHashTableWithCollisionList{}

	err = decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)

	return output, nil
}
