package system

import (
	"SDCC-Project-WorkerNode/utility"
	"net"
	"net/rpc"
)

const (
	DefaultNetwork                      = "tcp"
	AmazonAWSRegion                     = "us-east-1"
	AmazonS3BucketName                  = "graziani-filestorage"
	DefaultArbitraryFaultToleranceLevel = 3
)

func StartAcceptingRPCRequest(serviceTypeRequest interface{}, address string) error {

	var listener net.Listener
	var err error

	if listener, err = net.Listen(DefaultNetwork, address); err != nil {
		return err
	}

	if err = rpc.Register(serviceTypeRequest); err != nil {
		return err
	}

	defer func() {
		utility.CheckError(listener.Close())
	}()
	for {
		rpc.Accept(listener)
	}
}
