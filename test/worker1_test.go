package test

import (
	"fmt"
	"sync"
	"testing"
)

func Test_worker1(t *testing.T) {
	//heartbeat.SendHeartBeatToLeader(1)

	slice := make([]int, 20)
	var myWaitGroup sync.WaitGroup

	myWaitGroup.Add(20)

	for i := 0; i < len(slice); i++ {

		go test2(&slice[i], &myWaitGroup)
		/*
			go func(index int) {
				channel := make(chan int, 1)
				go test(channel)
				slice[index]=<-channel
				myWaitGroup.Done()
			}(i)*/
	}

	myWaitGroup.Wait()
	fmt.Print(slice)

}

func test2(p *int, group *sync.WaitGroup) {
	(*p) = 4
	(*group).Done()
}

func test(output chan int) {
	output <- 3
}

type Sssss struct {
	ff int
	tt int
}

func test3(fff chan Sssss) {

}
