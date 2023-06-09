package handlers

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	d "github.com/vtpaulsen/blockchain/v2/dataStructs"
	m "github.com/vtpaulsen/blockchain/v2/messages"
	p "github.com/vtpaulsen/blockchain/v2/peer"
	r "github.com/vtpaulsen/blockchain/v2/rsa"
)

func HandleConnection(conn net.Conn, peer *p.Peer) {
	fmt.Print("Handling connection")
	for {
		message := d.Message{TransactionType: "", Payload: new(interface{})}
		dec := p.GetDecoder(peer, conn)
		err := dec.Decode(&message)
		if err != nil {
			fmt.Println("handleConnection() error, " + err.Error())
			break
		}
		fmt.Println("Received message: " + message.TransactionType)

		switch message.TransactionType {
		case d.JOIN_MSG:
			requestMessage := d.Message{d.REQUEST_MSG, ""}
			m.SendMessage(requestMessage, peer, conn)
		case d.REQUEST_MSG:
			msgStruct := d.Message{d.REQUEST_REPLY, "127.0.0.1 " + strconv.Itoa(p.GetPort(peer))}
			m.SendMessage(msgStruct, peer, conn)
		case d.REQUEST_REPLY:
			msgStruct := d.Message{d.ADD_MSG, message.Payload.(string)}
			m.FloodMessage(msgStruct, peer)
		case d.ADD_MSG:
			checkConnection(message.Payload.(string), peer)
		case d.TRANSACTION:
			transaction := message.Payload.(d.TransactionStruct)
			p.Transaction(peer, transaction.From, transaction.To, transaction.Amount)
		case d.SIGNED_TRANSACTION:

			signedTransaction := message.Payload.(d.SignedTransactionStruct)

			from_pk := strings.Split(signedTransaction.From, ",")
			a := []string{signedTransaction.ID, from_pk[0], from_pk[1]}

			if r.Verify(a, signedTransaction.Signature) {
				p.SignedTransaction(peer, signedTransaction.From, signedTransaction.To, signedTransaction.Amount)
			} else {
				fmt.Println("The signature is not valid")
			}
		case d.PRINTLEDGER:
			p.PrintLedgers(peer)
		case d.MSG:
			fmt.Println(message)
		default:
			fmt.Println("Message does not make sense")
		}
	}
}

func checkConnection(msg string, peer *p.Peer) {
	splitMsg := strings.Split(msg, " ")
	port, err := strconv.Atoi(splitMsg[1])
	if err != nil {
		fmt.Println("Error parsing port")
		return
	}
	if p.GetPort(peer) != port {
		Connect(peer, splitMsg[0], port)
	} else {
		fmt.Println("Connection already exists")
	}
}
