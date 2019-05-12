package primarynode

import (
	"fmt"
	"sync"
	"testing"
)

func TestLeaderElection(t *testing.T) {

	var waitGroup sync.WaitGroup

	primaryNode1 := New(1)
	primaryNode2 := New(2)
	primaryNode3 := New(3)
	primaryNode4 := New(4)

	waitGroup.Add(1)

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
