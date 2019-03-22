package services

import (
	"SDCC-Project-WorkerNode/src/core/utility"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func startReceivingHeartBeating() {

	var socket *net.UDPConn

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 5000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := socket.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
	}
}

func startClientHeartBeating(primaryAddress string) {

	var addressOfUDPEndPoint *net.UDPAddr
	var socket *net.UDPConn
	var err error

	if addressOfUDPEndPoint, err = net.ResolveUDPAddr("udp", primaryAddress); err != nil {
		panic(err)
	}

	if socket, err = net.DialUDP("udp", nil, addressOfUDPEndPoint); err != nil {
		panic(err)
	}
	defer func() {
		utility.CheckError(socket.Close())
	}()

	done := make(chan bool, 1)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func(done chan bool) {
		for t := range ticker.C {
			fmt.Println("Sending beat at:", t)
			sendHeartBeat(socket)
		}
	}(done)

	<-done
}

func sendHeartBeat(socket *net.UDPConn) {
	enc := gob.NewEncoder(socket)
	enc.Encode("Im here")
}

func sendUDP(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
