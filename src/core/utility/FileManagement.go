package utility

import (
	"errors"
	"os"
)

func WriteDataToFile(data interface{}) error {
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
