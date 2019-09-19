package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	timeout = 10 * time.Minute
)

type Registry struct {
	dataTable  map[string][]byte
	timerTable map[string]*time.Timer
	path       string
}

func New(processID int, processType string, isRegisterVolatile bool) (*Registry, error) {

	output := new(Registry)

	(*output).timerTable = make(map[string]*time.Timer)
	(*output).dataTable = make(map[string][]byte)

	if !isRegisterVolatile {

		(*output).path = fmt.Sprintf("./registry-%s-%d/", processType, processID)

		if err := os.MkdirAll((*output).path, 0755); err != nil {
			return nil, err
		} else {
			return output, (*output).recoverDataFromDisk()
		}

	} else {
		(*output).path = ""
	}

	return output, nil
}

func (obj *Registry) Get(guid string) []byte {
	return (*obj).dataTable[guid]
}

func (obj *Registry) Set(key string, input []byte) error {

	if (*obj).dataTable[key] == nil {

		(*obj).dataTable[key] = input
		(*obj).timerTable[key] = time.NewTimer(timeout)

		if (*obj).path != "" {

			if err := (*obj).writeDataOnDisk(key, input); err != nil {
				panic(err)
			}
		}

		go (*obj).startAutomaticCleanerRoutine(key)

	} else {
		(*obj).timerTable[key].Reset(timeout)
	}

	return nil
}

func (obj *Registry) writeDataOnDisk(guid string, data []byte) error {

	if file, err := os.OpenFile(fmt.Sprintf("%s/%s", (*obj).path, guid), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		return err
	} else {

		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()

		if _, err := file.Write(data); err != nil {
			return err
		} else {
			return file.Sync()
		}
	}
}

func (obj *Registry) recoverDataFromDisk() error {

	if directory, err := ioutil.ReadDir((*obj).path); err != nil {
		return err
	} else {

		for _, entry := range directory {

			if !entry.IsDir() {

				if data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", (*obj).path, entry.Name())); err != nil {
					return err
				} else {

					(*obj).dataTable[entry.Name()] = data
					(*obj).timerTable[entry.Name()] = time.NewTimer(timeout)

					go (*obj).startAutomaticCleanerRoutine(entry.Name())

				}
			}
		}
		return nil
	}
}

func (obj *Registry) startAutomaticCleanerRoutine(guid string) {

	<-(*obj).timerTable[guid].C

	if (*obj).path != "" {

		filePath := fmt.Sprintf("%s/%s", (*obj).path, guid)
		if err := os.Remove(filePath); err != nil {
			panic(err)
		}
	}

	delete((*obj).dataTable, guid)
	delete((*obj).timerTable, guid)
}
