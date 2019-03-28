package primarynode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/utility"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	leaderNotElected = -1
)

type PrimaryNode struct {
	id                      int
	leaderId                int
	listenPortForRPC        string
	listenUDPPort           int
	allPrimaryNodeAddresses []string
}

func New(primaryNodeId int, listenUDPPort int) *PrimaryNode {

	output := new(PrimaryNode)

	(*output).id = primaryNodeId
	(*output).leaderId = leaderNotElected
	(*output).listenUDPPort = listenUDPPort
	(*output).allPrimaryNodeAddresses = []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003", "localhost:5004"}

	return output
}

func (obj *PrimaryNode) startToRespondToRPCRequests() {

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Request{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()
}

func sendingHeartbeat(destinationAddress string) {

	var addressOfUDPEndPoint *net.UDPAddr
	var socket *net.UDPConn
	var err error

	addressOfUDPEndPoint, err = net.ResolveUDPAddr("udp", destinationAddress)
	utility.CheckError(err)

	socket, err = net.DialUDP("udp", nil, addressOfUDPEndPoint)
	utility.CheckError(err)
	defer func() {
		utility.CheckError(socket.Close())
	}()

	//ticker := time.NewTicker(500 * time.Millisecond)

	for range time.NewTicker(500 * time.Millisecond).C {

		fmt.Println("Sending beat to ", addressOfUDPEndPoint.IP, ":", addressOfUDPEndPoint.Port)

		enc := gob.NewEncoder(socket)
		err = enc.Encode("Im here")
		utility.CheckError(err)
	}
}

func (obj *PrimaryNode) startSendingHeartbeatToBackups() {

	for index := 0; index < len((*obj).allPrimaryNodeAddresses); index++ {

		if index != (*obj).id {
			sendingHeartbeat((*obj).allPrimaryNodeAddresses[index])
		}
	}
}

func (obj *PrimaryNode) startReceivingHeartbeatFromLeader(listenUDPPort int, leaderNotRespondingChannel chan bool) {

	var socket *net.UDPConn
	var timerChannel chan time.Time

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: listenUDPPort,
		IP:   net.ParseIP("localhost"),
	}

	socket, err := net.ListenUDP("udp", &addr)
	utility.CheckError(err)
	defer func() {
		utility.CheckError(socket.Close())
	}()

	time.After(time.Second * 2)

	for {
		select {

		case <-timerChannel:

			fmt.Println("Expired")
			leaderNotRespondingChannel <- true
			return

		default:
			_, _, err = socket.ReadFromUDP(p)
			utility.CheckError(err)
			fmt.Println(p)
			//timer.Reset(time.Second)
		}
	}
}

func (obj *PrimaryNode) StartWork() {

	var leaderNotRespondingChannel chan bool

	for {

		if (*obj).leaderId == leaderNotElected {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose any leader!")
			(*obj).leaderId = zookeeper.StartLeaderElection((*obj).id)
			fmt.Println("Actual leader is ", (*obj).leaderId)

		} else if (*obj).leaderId != (*obj).id {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose leader id ", (*obj).leaderId)

			go (*obj).startReceivingHeartbeatFromLeader((*obj).listenUDPPort, leaderNotRespondingChannel)

			<-leaderNotRespondingChannel

			fmt.Println("I'm node id ", (*obj).id, ": leader doesn't respond")

			(*obj).leaderId = leaderNotElected

		} else {

			//go (*obj).startSendingHeartbeatToBackups()

			time.Sleep(time.Second * 5)

		}
	}
}
