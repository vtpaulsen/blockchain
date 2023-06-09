package handlers

import (
	"fmt"
	"net"
	"strconv"

	p "github.com/vtpaulsen/blockchain/v2/peer"
)

func Connect(peer *p.Peer, serverIp string, serverPort int) {
	serverPortString := strconv.Itoa(serverPort)
	conn, err := net.Dial("tcp", serverIp+":"+serverPortString)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connected to " + serverIp + ":" + serverPortString)

		p.AddToConnections(peer, conn)
		go HandleConnection(conn, peer)
	}
}
