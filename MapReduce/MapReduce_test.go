package MapReduce

import (
	"SDCC-Project/MapReduce/Node/Primary"
	"SDCC-Project/MapReduce/Node/Worker"
	"testing"
)

func Test_primary(t *testing.T) {

	worker1 := Worker.New(1, "127.0.0.1:12001", "127.0.0.1:13001")
	worker2 := Worker.New(2, "127.0.0.1:12002", "127.0.0.1:13002")
	worker3 := Worker.New(3, "127.0.0.1:12003", "127.0.0.1:13003")
	worker4 := Worker.New(4, "127.0.0.1:12004", "127.0.0.1:13004")
	worker5 := Worker.New(5, "127.0.0.1:12005", "127.0.0.1:14005")
	worker6 := Worker.New(6, "127.0.0.1:12006", "127.0.0.1:15006")

	go worker1.StartWork()
	go worker2.StartWork()
	go worker3.StartWork()
	go worker4.StartWork()
	go worker5.StartWork()
	go worker6.StartWork()

	Primary.New(1, "127.0.0.1:15001").StartWork()
}
