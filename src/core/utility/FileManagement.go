package utility

import (
	"encoding/gob"
	"errors"
	"os"
)

func WriteToLocalDisk(data interface{}) error {

	var outputFile *os.File
	var outputFileName string
	var err error

	if outputFileName, err = SHA512(data); err != nil {
		return err
	}
	if outputFile, err = os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		return err
	}
	defer func() {
		CheckError(outputFile.Close())
	}()

	encoder := gob.NewEncoder(outputFile)

	if err = encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func ReadFirst256ByteFromFile(filePath string) ([]byte, error) {

	if filePath == "" {
		return nil, errors.New(InvalidInput)
	}
	inputFile, errorOpeningFile := os.OpenFile(filePath, os.O_RDONLY, 0666)

	if errorOpeningFile != nil {
		return nil, errorOpeningFile
	}
	defer func() {
		CheckError(inputFile.Close())
	}()

	output := make([]byte, 256)

	if _, errorWhileReading := inputFile.Read(output); errorWhileReading != nil {
		return nil, errorWhileReading
	}

	return output, nil
}
