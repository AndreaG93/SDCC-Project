package worker

import "testing"

func Test_worker1(t *testing.T) {
	New(1, "127.0.0.1").StartWork()
}
