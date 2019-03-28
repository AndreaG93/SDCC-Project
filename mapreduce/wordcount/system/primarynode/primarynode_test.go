package primarynode

import (
	"fmt"
	"sync"
	"testing"
)

func TestLeaderElection(t *testing.T) {

	var waitGroup sync.WaitGroup

	primaryNode0 := New(0, 5000)
	primaryNode1 := New(1, 5001)
	primaryNode2 := New(2, 5002)
	primaryNode3 := New(3, 5003)
	primaryNode4 := New(4, 5004)

	waitGroup.Add(1)

	go func() {
		(*primaryNode0).StartWork()
		fmt.Println("PrimaryNode0 returned.")
	}()

	go func() {
		(*primaryNode1).StartWork()
		fmt.Println("PrimaryNode1 returned.")
	}()

	go func() {
		(*primaryNode2).StartWork()
		fmt.Println("PrimaryNode2 returned.")
	}()

	go func() {
		(*primaryNode3).StartWork()
		fmt.Println("PrimaryNode3 returned.")
	}()

	go func() {
		(*primaryNode4).StartWork()
		fmt.Println("PrimaryNode4 returned.")
	}()

	waitGroup.Wait()
}

func Test_GeneralSystemTest(t *testing.T) {

	/*
		go func() {
			StartPrimaryServices("localhost:10000")
		}()

		go func() {
			StartWorkerServices("localhost:5000")
		}()

		go func() {
			StartWorkerServices("localhost:5001")
		}()

		go func() {
			StartWorkerServices("localhost:5002")
		}()

		go func() {
			StartWorkerServices("localhost:5003")
		}()

		time.Sleep(100 * time.Millisecond)

		wordcount.SendRequest("../../../test-input-data/input.txt", "localhost:10000")
	*/
}
