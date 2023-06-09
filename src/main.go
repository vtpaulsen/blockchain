package main

import (
	"bufio"

	"fmt"
	"os"
	"strconv"
	"strings"

	h "github.com/vtpaulsen/blockchain/v2/handlers"
	m "github.com/vtpaulsen/blockchain/v2/messages"
	p "github.com/vtpaulsen/blockchain/v2/peer"
)

func main() {

	peer := p.InitializePeer()
	ip, port := promptUser()

	// Phase 1: Connecting to a server
	h.Connect(&peer, ip, port)

	// Phase 2: Making ourself a server
	go h.MakeServer(&peer)

	// Phase 3: Manuel input
	manuelInput(&peer)
}

func manuelInput(peer *p.Peer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		if m.CheckMessage(msg, peer) {
			message := m.ConvertMsgToStruct(msg)
			m.FloodMessage(message, peer)
		}
	}
}

func promptUser() (string, int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What's the IP address?")
	ip, ip_error := reader.ReadString('\n')
	fmt.Println("What's the port number?")
	port, port_error := reader.ReadString('\n')

	if ip_error != nil || port_error != nil {
		panic("The IP or port is invalid.")
	}

	trimIP := strings.TrimSpace(ip)
	if port != "\n" {
		portInt, err := strconv.Atoi(strings.TrimSpace(port))
		if err != nil {
			panic("The port has to be of type int.")
		}
		return trimIP, portInt
	}
	return trimIP, 0
}

/*
Tjek først om personen findes
Hvis personen findes, så tjekker vi om der er penge nok
Hvis der er penge nok, så trækker vi penge fra kontoen og returnere true
*/

/*
Ting vi ved der skal fikses:
1. Vores manualInput låser når vi prøver at withdraw
2. Vi mangler at give ledger med når en ny peer joiner så den får en opdateret kontoliste
*/
