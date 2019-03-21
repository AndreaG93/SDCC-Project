package utility

import (
	"core"
	"net"
	"net/rpc"
)

func StartRemoteService(serviceType interface{}, servicePort uint) error {

	var listener net.Listener
	var err error

	if listener, err = net.Listen(core.DefaultNetwork, "localhost:"+string(servicePort)); err != nil {
		return err
	}

	if err = rpc.Register(serviceType); err != nil {
		return err
	}

	defer func() {
		CheckError(listener.Close())
	}()
	for {
		rpc.Accept(listener)
	}
}
