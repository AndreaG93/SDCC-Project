package utility

import (
	"testing"
)

func Test_SerializationDeserialization(t *testing.T) {

	done := make(chan bool, 1)

	go func() {
		startClientHeartBeating("127.0.0.1:5000")
	}()

	go func() {
		startReceivingHeartBeating()
	}()

	<-done
}
