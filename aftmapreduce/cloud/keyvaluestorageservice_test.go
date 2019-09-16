package cloud

import (
	"SDCC-Project/aftmapreduce/cloud/amazons3"
	"SDCC-Project/aftmapreduce/utility"
	"io/ioutil"
	"testing"
)

const (
	key           = "test"
	testFile1Path = "./test1.txt"
)

func Test_BasicOperations(t *testing.T) {

	var err error
	var data1 []byte
	var data2 []byte

	keyValueStorageService := amazons3.New()

	data1, err = ioutil.ReadFile(testFile1Path)
	utility.CheckError(err)

	err = (*keyValueStorageService).Put(key, data1)
	utility.CheckError(err)

	data2, err = (*keyValueStorageService).Get(key)
	utility.CheckError(err)

	if !utility.Equal(data1, data2) {
		panic("Error")
	}

	utility.CheckError(keyValueStorageService.Remove(key))
}

func Test_BasicOperationsWithURL(t *testing.T) {

	var err error
	var data1 []byte
	var data2 []byte

	keyValueStorageService := amazons3.New()

	data1, err = ioutil.ReadFile(testFile1Path)
	utility.CheckError(err)

	err = (*keyValueStorageService).Put(key, data1)
	utility.CheckError(err)

	url, err := keyValueStorageService.RetrieveURLForGetOperation(key)
	utility.CheckError(err)

	data2, err = Download(url)
	utility.CheckError(err)

	if !utility.Equal(data1, data2) {
		panic("Error")
	}

	utility.CheckError(keyValueStorageService.Remove(key))
}
