package aftmapreduce

import (
	"SDCC-Project/utility"
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

	var listener net.Listener

	listener, _ = net.Listen("tcp", address)
	rpc.Register(serviceTypeRequest)

	defer func() {
		utility.CheckError(listener.Close())
	}()
	for {
		rpc.Accept(listener)
	}
}
