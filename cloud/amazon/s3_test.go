package amazon

import (
	"SDCC-Project/utility"
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

	(*amazonS3Client).Upload(testFile1Path, testKeyName)

	outputFile, err := os.Create(testFile2Path)
	utility.CheckError(err)

	(*amazonS3Client).Download(testKeyName, outputFile)

	digest1, _ := utility.GenerateDigestOfFileUsingSHA512(testFile1Path)
	digest2, _ := utility.GenerateDigestOfFileUsingSHA512(testFile2Path)

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

	(*amazonS3Client).Upload(testFile1Path, testKeyName)

	outputFile, err := ioutil.TempFile("", "testFile2Path")
	utility.CheckError(err)

	defer func() {
		utility.CheckError(os.Remove(outputFile.Name()))
	}()

	(*amazonS3Client).Download(testKeyName, outputFile)

	fileInfo, err := outputFile.Stat()
	buffer := make([]byte, fileInfo.Size())

	outputFile.Read(buffer)

	fmt.Println(string(buffer))

	(*amazonS3Client).Delete(testKeyName)
}
