package BFTMapReduce

import (
	"SDCC-Project/utility"
	"net"
	"net/rpc"
)

const (
	MapReduceRPCBasePort    = 20000
	MapReduceGetRPCBasePort = 30000

	DefaultNetwork     = "tcp"
	AmazonAWSRegion    = "us-east-1"
	AmazonS3BucketName = "graziani-filestorage"

	DefaultArbitraryFaultToleranceLevel = 3

	RPCPort = 30000
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
