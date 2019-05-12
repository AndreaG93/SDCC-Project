package heartbeat

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/nodestatusregister"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/timerregister"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"net"
	"strconv"
	"syscall"
	"time"
)

const (
	heartBeatsPort = 20000
)

func SendStoppableHeartBeatsTo(myNodeId uint, recipientNodeId uint, stop *bool) {

	var dataToSend []byte
	var socket *net.UDPConn
	var err error

	//ip, err := net.LookupIP("primary" + strconv.Itoa(recipientNodeId))
	ip, err := net.LookupIP("localhost")
	utility.CheckError(err)

	recipientUDPAddress, err := net.ResolveUDPAddr("udp", ip[0].String()+":"+strconv.Itoa(heartBeatsPort+int(recipientNodeId)))
	utility.CheckError(err)

	for {
		socket, err = net.DialUDP("udp", nil, recipientUDPAddress)
		if err != nil {
			break
		}
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	dataToSend, _ = utility.Encode(myNodeId)

	for _ = range ticker.C {

		if *stop {

			err = socket.Close()
			utility.CheckError(err)

			return
		}

		_, err = socket.Write(dataToSend)
		utility.CheckError(err)
	}

}

func SendHeartBeatsTo(myNodeId uint, recipientNodeId uint) {

	alwaysFalse := false
	SendStoppableHeartBeatsTo(myNodeId, recipientNodeId, &alwaysFalse)
}

func ReceiveHeartBeats(myNodeId uint, leaderId *uint) {

	var dataToReceive []byte
	var err error
	var socket *net.UDPConn
	var nodeID uint

	nodeStatusRegister := nodestatusregister.GetInstance()
	timerRegister := timerregister.GetInstance()

	dataToReceive = make([]byte, 4)

	socket, err = net.ListenUDP("udp", buildListenUDPAddress(myNodeId))
	utility.CheckError(err)

	defer func() {
		err = socket.Close()
		utility.CheckError(err)
	}()

	for {
		_, _, err = socket.ReadFromUDP(dataToReceive)
		utility.CheckError(err)

		err = utility.Decode(dataToReceive, &nodeID)
		utility.CheckError(err)

		if nodeStatusRegister.IsNodeOffline(nodeID) {

			nodeStatusRegister.SetNodeStatusAsOnline(nodeID)

			go func() {
				timerRegister.StartTimer(nodeID)

				fmt.Printf("Timer associated to %d expired", nodeID)
				nodeStatusRegister.SetNodeStatusAsOffline(nodeID)

				if (*leaderId) == nodeID {

					err = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
					utility.CheckError(err)
				}

			}()

		} else {
			timerRegister.StopResetAndRestart(nodeID)
		}
	}
}

func buildListenUDPAddress(myNodeId uint) *net.UDPAddr {

	listeningAddress := new(net.UDPAddr)

	ip, err := net.LookupIP("localhost")
	utility.CheckError(err)

	(*listeningAddress).Port = heartBeatsPort + int(myNodeId)
	(*listeningAddress).IP = net.ParseIP(ip[0].String())

	return listeningAddress
}
