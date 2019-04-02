package heartbeat

import (
	"sync"
	"testing"
)

func TestHeartbeatService(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)

	go StartToReceiveHeartbeat()

	StartToSendHeartbeat("Worker 1", "localhost:7000/heartbeat")
	StartToSendHeartbeat("Worker 2", "localhost:7000/heartbeat")
	StartToSendHeartbeat("Worker 3", "localhost:7000/heartbeat")

	wg.Wait()
}
