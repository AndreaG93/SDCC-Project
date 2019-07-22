package WorkerMutex

import (
	"sync"
)

var mutex *sync.Mutex
var once sync.Once

func GetInstance() *sync.Mutex {
	once.Do(func() {
		mutex = &sync.Mutex{}
	})
	return mutex
}
