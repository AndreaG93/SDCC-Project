package utility

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
)

func WriteToLocalDisk(filename string, data interface{}) error {

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(data); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func ReadLocalFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func WriteToLocalDisk2(data interface{}) error {

	var outputFile *os.File
	var outputFileName string
	var err error

	if outputFileName, err = GenerateDigestOfDataUsingSHA512(data); err != nil {
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
