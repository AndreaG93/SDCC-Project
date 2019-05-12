package timerregister

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/nodeidsregister"
	"sync"
	"time"
)

const (
	heartBeatingFrequency = 1 * time.Second
)

var register *TimerRegister
var once sync.Once

type TimerRegister struct {
	value map[uint]*time.Timer
}

func GetInstance() *TimerRegister {
	once.Do(func() {
		register = build()
	})
	return register
}

func build() *TimerRegister {

	nodeIDsRegister := nodeidsregister.GetInstance()

	output := new(TimerRegister)

	(*output).value = make(map[uint]*time.Timer)

	for id := range nodeIDsRegister.GetNodeIDs() {
		(*output).value[uint(id)] = time.NewTimer(heartBeatingFrequency)
	}

	return output
}

func (obj *TimerRegister) StopResetAndRestart(id uint) {

	timer := (*obj).value[id]

	timer.Stop()
	timer.Reset(heartBeatingFrequency)
}

func (obj *TimerRegister) StartTimer(id uint) {

	timer := (*obj).value[id]

	<-timer.C
}
