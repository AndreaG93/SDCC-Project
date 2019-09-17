package wordcount

import (
	"time"
)

const timeout = 2 * time.Second

type AFTTask interface {
	GetReplyChannel() chan interface{}
	GetFaultToleranceLevel() int
	GetAvailableWorkerProcessesRPCInternetAddresses() []string
	DoWeHaveEnoughMatchingReplyAfter(lastReply interface{}) bool
	ExecuteRPCCallTo(fullRPCInternetAddress string)
	GetOutput() interface{}
}

func Execute(task AFTTask) interface{} {

	workerProcessesRPCInternetAddresses := task.GetAvailableWorkerProcessesRPCInternetAddresses()
	repliesReceived := 0
	RPCCallSent := 0
	timer := time.NewTimer(timeout)

	for ; RPCCallSent <= task.GetFaultToleranceLevel(); RPCCallSent++ {
		fullRPCInternetAddress := workerProcessesRPCInternetAddresses[RPCCallSent]
		task.ExecuteRPCCallTo(fullRPCInternetAddress)
	}

	for {
		select {
		case <-timer.C:

			if thereAreOtherAvailableWorkerProcesses(task, RPCCallSent) {
				fullRPCInternetAddress := workerProcessesRPCInternetAddresses[RPCCallSent]
				task.ExecuteRPCCallTo(fullRPCInternetAddress)
				RPCCallSent++
			} else {
				panic("not enough worker processes")
			}

		case reply := <-task.GetReplyChannel():

			timer.Stop()
			repliesReceived++

			if task.DoWeHaveEnoughMatchingReplyAfter(reply) {

				close(task.GetReplyChannel())
				return task.GetOutput()

			} else {

				if repliesReceived == RPCCallSent && thereAreOtherAvailableWorkerProcesses(task, RPCCallSent) {
					fullRPCInternetAddress := workerProcessesRPCInternetAddresses[RPCCallSent]
					task.ExecuteRPCCallTo(fullRPCInternetAddress)
					RPCCallSent++
				} else {
					panic("not enough worker processes")
				}

				timer.Reset(timeout)
			}
		}
	}
}

func thereAreOtherAvailableWorkerProcesses(task AFTTask, RPCCallSentSoFar int) bool {
	return RPCCallSentSoFar < len(task.GetAvailableWorkerProcessesRPCInternetAddresses())
}
