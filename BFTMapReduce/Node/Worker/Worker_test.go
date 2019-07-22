package Worker

import (
	"testing"
)

func Test_Worker(t *testing.T) {

	worker1 := New(1, "127.0.0.1")
	worker2 := New(2, "127.0.0.1")
	worker3 := New(3, "127.0.0.1")
	worker4 := New(4, "127.0.0.1")
	worker5 := New(5, "127.0.0.1")
	worker6 := New(6, "127.0.0.1")

	go worker1.StartWork()
	go worker2.StartWork()
	go worker3.StartWork()
	go worker4.StartWork()
	go worker5.StartWork()
	worker6.StartWork()
}
