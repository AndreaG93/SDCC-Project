package amazons3

import (
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	testKeyName   = "test"
	testFile1Path = "./test1.txt"
	testFile2Path = "./test2.txt"
)

func Test_AmazonS3BasicOperations(t *testing.T) {

	amazonS3Client := New()

	file, err := os.Open(testFile1Path)
	utility.CheckError(err)

	(*amazonS3Client).Upload(file, testKeyName)
	utility.CheckError(file.Close())

	outputFile, err := os.Create(testFile2Path)
	utility.CheckError(err)

	(*amazonS3Client).Download(testKeyName, outputFile)
	utility.CheckError(outputFile.Close())

	data1, err := ioutil.ReadFile(testFile1Path)
	utility.CheckError(err)

	data2, err := ioutil.ReadFile(testFile2Path)
	utility.CheckError(err)

	digest1 := utility.GenerateDigestUsingSHA512(data1)
	digest2 := utility.GenerateDigestUsingSHA512(data2)

	if strings.Compare(digest1, digest2) != 0 {
		panic("Error")
	}

	(*amazonS3Client).Delete(testKeyName)

	if err := os.Remove(testFile2Path); err != nil {
		panic(err)
	}
}

func Test_AmazonS3BasicOperationsWithTemporaryFile(t *testing.T) {

	amazonS3Client := New()

	file, err := os.Open("./test1.txt")
	utility.CheckError(err)
	defer func() {
		utility.CheckError(file.Close())
	}()

	(*amazonS3Client).Upload(file, testKeyName)

	outputFile, err := ioutil.TempFile("", "testFile2Path")
	utility.CheckError(err)

	defer func() {
		utility.CheckError(outputFile.Close())
		utility.CheckError(os.Remove(outputFile.Name()))
	}()

	(*amazonS3Client).Download(testKeyName, outputFile)

	fileInfo, err := outputFile.Stat()
	utility.CheckError(err)
	buffer := make([]byte, fileInfo.Size())

	_, err = outputFile.Read(buffer)
	utility.CheckError(err)

	fmt.Println(string(buffer))

	(*amazonS3Client).Delete(testKeyName)
}
