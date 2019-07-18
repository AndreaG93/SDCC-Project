package DataStructures

import (
	"time"
)

const (
	timeout = 15 * time.Second
)

type AutoCleanerHashTable struct {
	hashTable map[string][]byte
}

func Build() *AutoCleanerHashTable {

	output := new(AutoCleanerHashTable)
	(*output).hashTable = make(map[string][]byte)

	return output
}

func (obj *AutoCleanerHashTable) Get(key string) []byte {
	return (*obj).hashTable[key]
}

func (obj *AutoCleanerHashTable) Set(key string, data []byte) {

	(*obj).hashTable[key] = data

	go (*obj).automaticClean(key)
}

func (obj *AutoCleanerHashTable) automaticClean(digest string) {

	timer1 := time.NewTimer(timeout)
	<-timer1.C

	delete((*obj).hashTable, digest)
}
