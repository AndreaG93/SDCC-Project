package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestHe(t *testing.T) {

	var stopChannel = make(chan os.Signal, 2)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func(stopChannel chan os.Signal) {
		if err := system.StartAcceptingRPCRequest(&wordcount.Request{}, "localhost:4000", stopChannel); err != nil {
			panic(err)
		}
	}(stopChannel)

	time.Sleep(3 * time.Second)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	time.Sleep(15 * time.Second)
}
