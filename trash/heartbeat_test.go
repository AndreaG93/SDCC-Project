package trash

import (
	"SDCC-Project-WorkerNode"
	"testing"
)

func Test_SerializationDeserialization(t *testing.T) {

	done := make(chan bool, 1)

	go func() {
		SDCC_Project_WorkerNode.startClientHeartBeating("127.0.0.1:5000")
	}()

	go func() {
		SDCC_Project_WorkerNode.startReceivingHeartBeating()
	}()

	<-done
}
