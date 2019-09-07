package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/utility"
	"io/ioutil"
	"os"
	"strings"
)

func downloadSourceDataFromCloud(sourceDataDigest string) *os.File {

	output, err := ioutil.TempFile(os.TempDir(), sourceDataDigest)
	utility.CheckError(err)

	node.GetAmazonS3Client().Download(sourceDataDigest, output)

	_, err = output.Seek(0, 0)
	utility.CheckError(err)

	return output
}

func getSplits(sourceDataDigest string, splitsCardinality int) []string {

	var fileInfo os.FileInfo
	var err error

	output := make([]string, splitsCardinality)

	inputFile := downloadSourceDataFromCloud(sourceDataDigest)
	fileInfo, err = inputFile.Stat()
	utility.CheckError(err)

	averageChunkSize := fileInfo.Size() / int64(splitsCardinality)
	readByte := make([]byte, 1)

	for index, currentStartByte, currentEndByte := int64(0), int64(0), averageChunkSize; ; {

		_, err = inputFile.Seek(currentEndByte, 0)
		utility.CheckError(err)

		_, err = inputFile.Read(readByte)
		utility.CheckError(err)

		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)

			_, err = inputFile.Seek(currentStartByte, 0)
			utility.CheckError(err)

			_, err = inputFile.Read(readData)
			utility.CheckError(err)

			output[index] = string(readData)
			index++

			currentStartByte = currentEndByte

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {

				currentEndByte += averageChunkSize

			} else {

				readData := make([]byte, fileInfo.Size()-currentStartByte)

				_, err = inputFile.Seek(currentStartByte, 0)
				utility.CheckError(err)

				_, err = inputFile.Read(readData)
				utility.CheckError(err)

				output[index] = string(readData)
				break
			}

		} else {
			currentEndByte++
		}
	}

	utility.CheckError(inputFile.Close())
	utility.CheckError(os.Remove(inputFile.Name()))
	return output
}
