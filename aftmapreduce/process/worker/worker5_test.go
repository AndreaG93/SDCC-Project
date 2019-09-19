package worker

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"testing"
)

func Test_worker5(t *testing.T) {

	utility.CheckError(process.Initialize(5, 1, process.WorkerProcessType, "127.0.0.1", []string{"127.0.0.1:2181"}))
	StartWork()
}
