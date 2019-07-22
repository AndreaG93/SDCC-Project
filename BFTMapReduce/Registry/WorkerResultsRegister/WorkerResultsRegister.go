package WorkerResultsRegister

import (
	"SDCC-Project/BFTMapReduce/DataStructures/AutoCleanerHashTable"
	"sync"
)

var data *AutoCleanerHashTable.AutoCleanerHashTable
var once sync.Once

func GetInstance() *AutoCleanerHashTable.AutoCleanerHashTable {
	once.Do(func() {
		data = AutoCleanerHashTable.Build()
	})
	return data
}
