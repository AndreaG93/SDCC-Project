package heartbeat

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/nodestatusregister"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/timerregister"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"net"
	"syscall"
	"time"
)

const (
	leaderHeartBeatListenPort = 6000
	workerHeartBeatListenPort = 5000
)

func LeaderMonitoring(leaderId uint) {

	var err error
	var listeningSocket *net.UDPConn
	var nodeID uint
	var nodeStatusRegister *nodestatusregister.NodeStatusRegister
	var timerRegister *timerregister.TimerRegister

	nodeStatusRegister = nodestatusregister.GetInstance()
	timerRegister = timerregister.GetInstance()

	p := make([]byte, 4)

	socketListenAddress := buildSocketListenAddress(leaderHeartBeatListenPort)

	listeningSocket, err = net.ListenUDP("udp", socketListenAddress)
	utility.CheckError(err)

	for {
		_, _, err = listeningSocket.ReadFromUDP(p)
		utility.CheckError(err)

		err = utility.Decode(p, &nodeID)
		utility.CheckError(err)

		if nodeID != leaderId {
			return
		} else {

			if nodeStatusRegister.IsNodeOffline(nodeID) {

				nodeStatusRegister.SetNodeStatusAsOnline(nodeID)

				go func() {
					timerRegister.StartTimer(nodeID)

					err = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
					utility.CheckError(err)

					nodeStatusRegister.SetNodeStatusAsOffline(nodeID)
				}()

			} else {
				timerRegister.StopResetAndRestart(nodeID)
			}
		}
	}
}

func WorkerMonitoring() {

	var err error
	var listeningSocket *net.UDPConn
	var nodeID uint
	var nodeStatusRegister *nodestatusregister.NodeStatusRegister
	var timerRegister *timerregister.TimerRegister

	nodeStatusRegister = nodestatusregister.GetInstance()
	timerRegister = timerregister.GetInstance()

	p := make([]byte, 4)

	socketListenAddress := buildSocketListenAddress(workerHeartBeatListenPort)

	listeningSocket, err = net.ListenUDP("udp", socketListenAddress)
	utility.CheckError(err)

	for {
		_, _, err = listeningSocket.ReadFromUDP(p)
		utility.CheckError(err)

		err = utility.Decode(p, &nodeID)
		utility.CheckError(err)

		if nodeStatusRegister.IsNodeOffline(nodeID) {

			nodeStatusRegister.SetNodeStatusAsOnline(nodeID)

			go func() {
				timerRegister.StartTimer(nodeID)

				fmt.Printf("Timer associated to %d expired", nodeID)
				nodeStatusRegister.SetNodeStatusAsOffline(nodeID)
			}()

		} else {
			timerRegister.StopResetAndRestart(nodeID)
		}
	}
}

func StartHeartBeating(id uint, sendingSocketAddress string) {

	var dataToSend []byte
	var socket *net.UDPConn
	var err error

	dataToSend, _ = utility.Encode(id)

	sendingUDPSocketAddress, err := net.ResolveUDPAddr("udp", sendingSocketAddress)
	utility.CheckError(err)

	socket, err = net.DialUDP("udp", nil, sendingUDPSocketAddress)
	utility.CheckError(err)

	defer func() {
		utility.CheckError(socket.Close())
	}()

	done := make(chan bool, 1)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func(done chan bool) {
		for _ = range ticker.C {

			_, err = socket.Write(dataToSend)
			utility.CheckError(err)
		}
	}(done)

	<-done
}

func buildSocketListenAddress(port int) *net.UDPAddr {

	listeningAddress := new(net.UDPAddr)

	(*listeningAddress).Port = port
	(*listeningAddress).IP = net.ParseIP("127.0.0.1")

	return listeningAddress
}
