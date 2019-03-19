package data_structures

import (
	"core/utility"
	"encoding/gob"
	"fmt"
	"os"
)

type TokenWordHashTableSerialized []TokenWordHashTableSerializationUnit

type TokenWordHashTableSerializationUnit struct {
	TableSizeOrIndex uint
	Word             string
	Occurrences      uint
}

func (obj *TokenWordHashTableSerialized) writeOnLocalDisk() error {

	var outputFile *os.File
	var outputFileName string
	var err error

	if outputFileName, err = utility.SHA512(*obj); err != nil {
		return err
	}
	if outputFile, err = os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(outputFile.Close())
	}()

	encoder := gob.NewEncoder(outputFile)

	if err = encoder.Encode(*obj); err != nil {
		return err
	}

	return nil
}

func writeTokenWordHashTableSerializationUnitsOnLocalDisk(data []TokenWordHashTableSerializationUnit) error {

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

func readSerializedWordHashTableWithCollisionListFromLocalDisk(filePath string) ([]TokenWordHashTableSerializationUnit, error) {

	var inputFile *os.File
	var err error

	if inputFile, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return nil, err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()

	decoder := gob.NewDecoder(inputFile)
	output := []TokenWordHashTableSerializationUnit{}

	err = decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)

	return output, nil
}
