package worker

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"testing"
)

func Test_worker2(t *testing.T) {
	utility.CheckError(process.Initialize(2, 0, process.WorkerProcessType, "localhost", []string{"127.0.0.1:2181"}))
	StartWork()
}
