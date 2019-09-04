package worker

import "testing"

func Test_worker0(t *testing.T) {
	Initialize(0, 0, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
