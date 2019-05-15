package test

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/heartbeat"
	"testing"
)

func Test_worker0(t *testing.T) {
	heartbeat.SendHeartBeatToLeader(0)
}
