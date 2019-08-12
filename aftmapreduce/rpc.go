package aftmapreduce

import (
	"SDCC-Project/utility"
	"net"
	"net/rpc"
)

const (
	WordCountMapTaskRPCBasePort       = 20000
	WordCountReduceTaskRPCBasePort    = 30000
	WordCountReceiveRPCBasePort       = 40000
	WordCountRequestRPCBasePort       = 50000
	WordCountSendRPCBasePort          = 60000
	WordCountDataRetrieverRPCBasePort = 70000
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
