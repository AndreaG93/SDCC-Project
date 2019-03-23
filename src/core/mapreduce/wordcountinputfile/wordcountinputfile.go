package wordcountinputfile

import (
	"SDCC-Project-WorkerNode/src/core/utility"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Split(inputFilename string, parts int64) error {

	var inputFile *os.File
	var fileInfo os.FileInfo
	var averageChunkSize int64
	var err error

	if inputFile, err = os.Open(inputFilename); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()
	if fileInfo, err = inputFile.Stat(); err != nil {
		return err
	}

	averageChunkSize = fileInfo.Size() / parts
	readByte := make([]byte, 1)

	for index, currentStartByte, currentEndByte := 0, int64(0), int64(averageChunkSize); ; {

		if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
			return err
		}
		if _, err = inputFile.Read(readByte); err != nil {
			return err
		}

		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)

			if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
				return err
			}
			if _, err = inputFile.Read(readData); err != nil {
				return err
			}
			if err := ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData, 0644); err != nil {
				return err
			}
			index++

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {

				currentStartByte = currentEndByte
				currentEndByte += averageChunkSize

			} else {

				readData := make([]byte, fileInfo.Size()-currentEndByte)

				if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
					return err
				}
				if _, err = inputFile.Read(readData); err != nil {
					return err
				}
				if err := ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData, 0644); err != nil {
					return err
				}
				break
			}

		} else {
			currentEndByte++
		}

	}
	return nil
}
