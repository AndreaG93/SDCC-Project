package worker

import "testing"

func Test_worker3(t *testing.T) {
	New(3, "127.0.0.1").StartWork()
}