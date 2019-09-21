package worker

import (
	"fmt"
	"testing"
)

const golangCompilerName = "go.exe"

func TestAbs(t *testing.T) {

	output := make(chan int)

	go func() {
		fmt.Printf("1 - %d\n", <-output)
	}()
	go func() {
		fmt.Printf("2 - %d\n", <-output)
	}()

	output <- 5
	output <- 5
}
