package heartbeat

import (
	"testing"
)

func Test_SerializationDeserialization(t *testing.T) {

	done := make(chan bool, 1)

	commits := map[uint]bool{
		uint(0): false,
		uint(1): false,
	}

	go func() {

		startWorkerNodeHeartBeating(0, "127.0.0.1:5000")

		startWorkerNodeHeartBeating(0, "127.0.0.1:5000")
	}()

	go func() {

		startWorkerNodeHeartBeating(1, "127.0.0.1:5000")

	}()

	go func() {
		StartWorkerNodesHeartBeatingMonitoring(commits)
	}()

	<-done
}
