package worker

import "testing"

func Test_worker2(t *testing.T) {
	New(2, "127.0.0.1", []string{"127.0.0.1"}).StartWork()
}
