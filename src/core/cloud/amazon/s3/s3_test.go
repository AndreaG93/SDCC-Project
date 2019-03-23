package s3

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_AmazonS3(t *testing.T) {

	const (
		contentOfFile = "test"
		fileName      = "./input.txt"
		objectKeyName = "TestKey"
	)

	var readDataFromDownloadedFile []byte
	var err error

	if err = UploadFileToCloudStorage(fileName, objectKeyName); err != nil {
		panic(err)
	}

	if err = os.Remove(fileName); err != nil {
		panic(err)
	}

	if err = DownloadFileFromCloudStorage(fileName, objectKeyName); err != nil {
		panic(err)
	}

	if readDataFromDownloadedFile, err = ioutil.ReadFile(fileName); err != nil {
		panic(err)
	}

	if strings.Compare(string(readDataFromDownloadedFile), contentOfFile) != 0 {
		panic(err)
	}

	if err = DeleteFileFromCloudStorage(objectKeyName); err != nil {
		panic(err)
	}
}
