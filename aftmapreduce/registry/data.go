package registry

import (
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	timeout = 30 * time.Second
)

type DataRegistry struct {
	dataRegistryFolderPath string
	content                map[string][]byte
	permanentRegistry      bool
	timerTable             map[string]*time.Timer
}

func NewDataRegistry(nodeName string, writeRegistryOnFileSystem bool) *DataRegistry {

	output := new(DataRegistry)

	(*output).timerTable = make(map[string]*time.Timer)
	(*output).permanentRegistry = writeRegistryOnFileSystem

	if (*output).permanentRegistry {
		(*output).dataRegistryFolderPath = fmt.Sprintf("./registry-%s/", nodeName)
		utility.CheckError(os.MkdirAll((*output).dataRegistryFolderPath, 0755))
		(*output).initializeDataRegistry((*output).dataRegistryFolderPath)

	} else {
		(*output).content = make(map[string][]byte)
	}

	return output
}

func (obj *DataRegistry) Get(key string) []byte {
	return (*obj).content[key]
}

func (obj *DataRegistry) Set(key string, input []byte) {

	if (*obj).content[key] == nil {

		(*obj).content[key] = input
		(*obj).timerTable[key] = time.NewTimer(timeout)

		if (*obj).permanentRegistry {
			(*obj).writeOnDisk(key, input)
		}

		go (*obj).automaticClean(key)

	} else {
		(*obj).timerTable[key].Reset(timeout)
	}
}

func (obj *DataRegistry) automaticClean(digest string) {

	<-(*obj).timerTable[digest].C

	if (*obj).permanentRegistry {

		path := fmt.Sprintf("%s%s", (*obj).dataRegistryFolderPath, digest)
		utility.CheckError(os.Remove(path))
	}

	delete((*obj).content, digest)
	delete((*obj).timerTable, digest)
}

func (obj *DataRegistry) writeOnDisk(digest string, data []byte) {

	file, err := os.OpenFile(fmt.Sprintf("./%s/%s", (*obj).dataRegistryFolderPath, digest), os.O_WRONLY|os.O_CREATE, 0666)
	utility.CheckError(err)
	defer func() {
		utility.CheckError(file.Close())
	}()

	_, err = file.Write(data)
	utility.CheckError(err)

	utility.CheckError(file.Sync())
}

func (obj *DataRegistry) initializeDataRegistry(dataRegistryDirectoryPath string) {

	(*obj).content = make(map[string][]byte)

	c, err := ioutil.ReadDir(dataRegistryDirectoryPath)
	utility.CheckError(err)

	for _, entry := range c {

		if !entry.IsDir() {
			rawData, err := ioutil.ReadFile(dataRegistryDirectoryPath + entry.Name())
			utility.CheckError(err)

			(*obj).content[entry.Name()] = rawData
			go (*obj).automaticClean(entry.Name())
		}
	}
}
