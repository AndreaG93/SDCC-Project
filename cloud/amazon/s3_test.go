package amazon

import (
	"SDCC-Project-WorkerNode/utility"
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

	(*amazonS3Client).Download(testKeyName, testFile2Path)

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
