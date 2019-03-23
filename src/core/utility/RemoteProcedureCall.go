package utility

import (
	"SDCC-Project-WorkerNode/src/core"
	"net"
	"net/rpc"
)

func StartAcceptingServiceRequest(serviceTypeRequest interface{}, serviceRequestListenPort uint) error {

	var listener net.Listener
	var err error

	if listener, err = net.Listen(core.DefaultNetwork, "localhost:"+string(serviceRequestListenPort)); err != nil {
		return err
	}

	if err = rpc.Register(serviceTypeRequest); err != nil {
		return err
	}

	defer func() {
		CheckError(listener.Close())
	}()
	for {
		rpc.Accept(listener)
	}
}
