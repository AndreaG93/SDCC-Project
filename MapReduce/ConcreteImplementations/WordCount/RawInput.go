package WordCount

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/utility"
	"os"
	"strings"
)

type RawInput struct {
	MapCardinality int64
	FileDigest     string
}

func (obj *RawInput) MapOutputRawDataToReduceInputData(mapOutputRawData [][]byte) []Data.Split {
	panic("implement me")
}

func (obj *RawInput) ReduceOutputRawDataToFinalOutput(ReduceOutputRawData [][]byte) []Data.Split {
	panic("implement me")
}

func (obj *RawInput) Split() []Data.Split {
	return nil
}

func (obj *RawInput) SplitInputFile() ([]string, error) {

	var inputFile *os.File
	var fileInfo os.FileInfo
	var err error

	output := make([]string, (*obj).MapCardinality)

	if inputFile, err = os.Open((*obj).FileDigest); err != nil {
		return nil, err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()
	if fileInfo, err = inputFile.Stat(); err != nil {
		return nil, err
	}

	averageChunkSize := fileInfo.Size() / (*obj).MapCardinality
	readByte := make([]byte, 1)

	for index, currentStartByte, currentEndByte := int64(0), int64(0), averageChunkSize; ; {

		if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
			return nil, err
		}
		if _, err = inputFile.Read(readByte); err != nil {
			return nil, err
		}

		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)

			if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
				return nil, err
			}
			if _, err = inputFile.Read(readData); err != nil {
				return nil, err
			}

			output[index] = string(readData)
			index++

			currentStartByte = currentEndByte

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {

				currentEndByte += averageChunkSize

			} else {

				readData := make([]byte, fileInfo.Size()-currentStartByte)

				if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
					return nil, err
				}
				if _, err = inputFile.Read(readData); err != nil {
					return nil, err
				}
				output[index] = string(readData)
				break
			}

		} else {
			currentEndByte++
		}
	}

	return output, nil
}
