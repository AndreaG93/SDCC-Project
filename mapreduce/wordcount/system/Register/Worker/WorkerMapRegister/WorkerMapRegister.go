package WorkerMapRegister

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/DataStructures"
	"sync"
)

var data *DataStructures.AutoCleanerHashTable
var once sync.Once

func GetInstance() *DataStructures.AutoCleanerHashTable {
	once.Do(func() {
		data = DataStructures.Build()
	})
	return data
}
