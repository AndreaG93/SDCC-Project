package worker

import "testing"

func Test_worker6(t *testing.T) {
	Initialize(5, 1, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}