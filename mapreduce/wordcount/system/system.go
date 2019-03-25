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

func StartAcceptingServiceRequest(serviceTypeRequest interface{}, serviceRequestListenPort uint) error {

	var listener net.Listener
	var err error

	if listener, err = net.Listen(DefaultNetwork, "localhost:"+string(serviceRequestListenPort)); err != nil {
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
