package worker

import "testing"

func Test_worker5(t *testing.T) {
	New(5, "127.0.0.1", []string{"127.0.0.1"}).StartWork()
}
