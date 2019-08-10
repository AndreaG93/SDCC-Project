package worker

import "testing"

func Test_worker5(t *testing.T) {
	Initialize(5, 2, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
