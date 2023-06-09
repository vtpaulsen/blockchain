package handlers

import (
	"fmt"
	"net"
	"strconv"

	d "github.com/vtpaulsen/blockchain/v2/dataStructs"
	m "github.com/vtpaulsen/blockchain/v2/messages"
	p "github.com/vtpaulsen/blockchain/v2/peer"
)

func MakeServer(peer *p.Peer) {
	ln, err := net.Listen("tcp", ":")
	printIPPort(peer, ln)
	if err != nil {
		fmt.Println("makeServer() error, " + err.Error())
	}
	JoinMsg := d.Message{d.JOIN_MSG, ""}
	m.FloodMessage(JoinMsg, peer)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("makeServer() error, " + err.Error())
		}
		p.AddToConnections(peer, conn)
		go HandleConnection(conn, peer)
	}
}

func printIPPort(peer *p.Peer, ln net.Listener) {
	ifaces, _ := net.Interfaces()
	var ip net.IP
	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}

	/* SKAL RYKKES */
	p.SetIP(peer, ip.String())
	p.SetPort(peer, ln.Addr().(*net.TCPAddr).Port)
	fmt.Println("Starting a server on:", ip.String())
	fmt.Println("Starting on IP: 127.0.0.1:" + strconv.Itoa(ln.Addr().(*net.TCPAddr).Port))
}
