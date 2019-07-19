package Worker

import (
	"testing"
)

func Test_primary2(t *testing.T) {

	worker1 := New(1, "127.0.0.1:12001", "127.0.0.1:13001")
	worker2 := New(2, "127.0.0.1:12002", "127.0.0.1:13002")
	worker3 := New(3, "127.0.0.1:12003", "127.0.0.1:13003")
	worker4 := New(4, "127.0.0.1:12004", "127.0.0.1:13004")
	worker5 := New(5, "127.0.0.1:12005", "127.0.0.1:14005")
	worker6 := New(6, "127.0.0.1:12006", "127.0.0.1:15006")

	go worker1.StartWork()
	go worker2.StartWork()
	go worker3.StartWork()
	go worker4.StartWork()
	go worker5.StartWork()
	worker6.StartWork()
}
