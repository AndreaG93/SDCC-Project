package worker

import "testing"

func Test_worker4(t *testing.T) {
	New(4, "127.0.0.1").StartWork()
}
