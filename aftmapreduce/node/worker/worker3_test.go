package worker

import "testing"

func Test_worker3(t *testing.T) {
	Initialize(3, 1, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
