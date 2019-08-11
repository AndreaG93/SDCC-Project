package worker

import "testing"

func Test_worker2(t *testing.T) {
	Initialize(1, 0, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
