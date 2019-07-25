package aftmapreduce

import (
	"SDCC-Project/utility"
	"net"
	"net/rpc"
)

const (
	MapReduceRPCBasePort        = 20000
	MapReduceGetRPCBasePort     = 30000
	MapReduceRequestRPCBasePort = 40000
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
