package test

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/heartbeat"
	"testing"
)

func Test_worker1(t *testing.T) {
	heartbeat.SendHeartBeatToLeader(1)
}
