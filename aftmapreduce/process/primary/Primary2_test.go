package primary

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"testing"
)

func Test_primary3(t *testing.T) {
	utility.CheckError(process.Initialize(2, 0, process.PrimaryProcessType, "127.0.0.1", []string{"127.0.0.1:2181"}))
	StartWork()
}
