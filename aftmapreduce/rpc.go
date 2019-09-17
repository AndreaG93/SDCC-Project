package aftmapreduce

import (
	"net"
	"net/rpc"
)

const (
	WordCountMapTaskRPCBasePort    = 2000
	WordCountReduceTaskRPCBasePort = 3000
	WordCountReceiveRPCBasePort    = 4000
	WordCountRequestRPCBasePort    = 5000
	WordCountSendRPCBasePort       = 6000
	WordCountRetrieverRPCBasePort  = 7000
)

func StartAcceptingRPCRequest(serviceTypeRequest interface{}, address string) {

	if listener, err := net.Listen("tcp", address); err != nil {
		panic(err)
	} else {

		if err := rpc.Register(serviceTypeRequest); err != nil {
			panic(err)
		} else {

			defer func() {
				if err := listener.Close(); err != nil {
					panic(err)
				}
			}()
			for {
				rpc.Accept(listener)
			}
		}
	}
}

func MakeRPCCall(serviceMethod string, internetAddress string, input interface{}, output interface{}, replyChannel chan interface{}) {

	if worker, err := rpc.Dial("tcp", internetAddress); err == nil {
		if err = worker.Call(serviceMethod, input, output); err == nil {
			replyChannel <- output.(interface{})
		}
	}
}
