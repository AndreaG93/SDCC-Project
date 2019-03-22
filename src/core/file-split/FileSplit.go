package file_split

import (
	"SDCC-Project-WorkerNode/src/core/utility"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func sss(inputFilename string, parts int64) error {

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

		inputFile.Seek(currentEndByte, 0)
		inputFile.Read(readByte)
		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)
			inputFile.Seek(currentStartByte, 0)
			inputFile.Read(readData)

			if err := ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData, 0644); err != nil {
				return err
			}
			index++

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {
				currentStartByte = currentEndByte
				currentEndByte += averageChunkSize
			} else {

				inputFile.Seek(currentEndByte, 0)
				readData := make([]byte, fileInfo.Size()-currentEndByte)

				inputFile.Read(readData)

				if err := ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData, 0644); err != nil {
					return err
				}

				break
			}

		} else {
			currentEndByte++
		}

	}

	inputFile.Read(readByte)

	fmt.Println(string(readByte[0]))
	fmt.Println(fileInfo.Size())
	fmt.Println(averageChunkSize)
	return nil
}

func SplitFile(inputFilename string, parts uint) error {

	var averageChunkSize uint
	var inputFileSize uint
	var readData []byte
	var err error

	if readData, err = ioutil.ReadFile(inputFilename); err != nil {
		return err
	}

	inputFileSize = uint(len(readData))
	averageChunkSize = inputFileSize / parts

	for index, currentStartByte, currentEndByte := 0, uint(0), averageChunkSize; ; {

		currentChar := string(readData[currentEndByte])

		if strings.Compare(currentChar, " ") == 0 {
			fmt.Println(inputFilename + string(index))
			err = ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData[currentStartByte:currentEndByte], 0644)
			if err != nil {
				return err
			}

			index++

			//fmt.Println(string(readData[currentStartByte:currentEndByte]))

			if (currentEndByte + averageChunkSize) < inputFileSize {
				currentStartByte = currentEndByte
				currentEndByte += averageChunkSize
			} else {

				ioutil.WriteFile(inputFilename+"_"+strconv.Itoa(index), readData[currentEndByte:], 0644)
				break
			}

		} else {
			currentEndByte++
		}
	}

	return nil
}
