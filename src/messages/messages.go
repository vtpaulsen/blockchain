package messages

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	d "github.com/vtpaulsen/blockchain/v2/dataStructs"
	p "github.com/vtpaulsen/blockchain/v2/peer"
	r "github.com/vtpaulsen/blockchain/v2/rsa"
)

func FloodMessage(msg d.Message, peer *p.Peer) {

	for conn := range p.GetConnections(peer) {
		fmt.Println("Sending message to " + conn.RemoteAddr().String() + ": " + msg.TransactionType)
		enc := p.GetEncoder(peer, conn)
		enc.Encode(msg)
	}
}

func SendMessage(msg d.Message, peer *p.Peer, conn net.Conn) {
	enc := p.GetEncoder(peer, conn)
	enc.Encode(msg)
}

func ConvertMsgToStruct(msg string) d.Message {
	splitMsg := strings.Split(msg, " ")
	ReturnMessage := d.Message{}
	switch splitMsg[0] {
	case d.SIGNED_TRANSACTION:
		amount, err := strconv.Atoi(strings.TrimSpace(splitMsg[5]))
		if err != nil {
			fmt.Println(err)
		}
		fromString := splitMsg[1] + "," + splitMsg[2]
		toString := splitMsg[3] + "," + splitMsg[4]
		Signed_Transaction := d.SignedTransactionStruct{
			ID:        splitMsg[0],
			From:      fromString,
			To:        toString,
			Amount:    amount,
			Signature: r.Sign(splitMsg[:6], splitMsg[6]),
		}
		return d.Message{TransactionType: d.SIGNED_TRANSACTION, Payload: Signed_Transaction}
	case d.TRANSACTION:
		amount, err := strconv.Atoi(strings.TrimSpace(splitMsg[3]))
		if err != nil {
			fmt.Println(err)
		}
		Transaction := d.TransactionStruct{
			ID:     splitMsg[0],
			From:   splitMsg[1],
			To:     splitMsg[2],
			Amount: amount,
		}
		return d.Message{TransactionType: d.TRANSACTION, Payload: Transaction}
	case d.PRINTLEDGER:
		return d.Message{TransactionType: d.PRINTLEDGER, Payload: nil}
	case d.MSG:
		return d.Message{TransactionType: d.MSG, Payload: msg}
	}
	return ReturnMessage
}

func CheckMessage(msg string, peer *p.Peer) bool {
	splitMsg := strings.Split(msg, " ")

	switch splitMsg[0] {
	case d.ADD_TO_LEDGER:
		if len(splitMsg) >= 2 {
			p.AddToLedgers(peer, splitMsg[1])
			return true
		}
	case d.WITHDRAW:
		if len(splitMsg) >= 3 {
			amount, _ := strconv.Atoi(splitMsg[2])
			if p.IsWithdrawLegal(peer, splitMsg[1], amount) {
				p.Withdraw(peer, splitMsg[1], amount)
				return true
			}
		}
	case d.DEPOSIT:
		if len(splitMsg) >= 3 {
			if p.DoesPersonExist(peer, splitMsg[1]) {
				amount, _ := strconv.Atoi(strings.TrimSpace(splitMsg[2]))
				p.Deposit(peer, splitMsg[1], amount)
				return true
			}
		}
	case d.TRANSACTION:
		if len(splitMsg) == 4 {
			amount, err := strconv.Atoi(strings.TrimSpace(splitMsg[3]))
			if err != nil {
				fmt.Println(err.Error())
			}
			if p.IsWithdrawLegal(peer, splitMsg[1], amount) {
				p.Transaction(peer, splitMsg[1], splitMsg[2], amount)
				return true
			} else {
				return false
			}
		}
	case d.SIGNED_TRANSACTION:
		if len(splitMsg) == 7 {
			amount, err := strconv.Atoi(strings.TrimSpace(splitMsg[5]))
			if err != nil {
				fmt.Println(err.Error())
			}

			sign := r.Sign(splitMsg[:3], splitMsg[6])

			if !r.Verify(splitMsg[:3], sign) {
				fmt.Println("The signature is not valid")
				return false
			}

			if p.IsWithdrawLegal(peer, splitMsg[1], amount) {
				p.SignedTransaction(peer, splitMsg[1], splitMsg[3], amount)
				return true
			} else {
				return false
			}
		}
		return true

	case d.PRINTLEDGER:
		p.PrintLedgers(peer)
		return true
	case d.MSG:
		return true
	case d.GET_ACCOUNT:
		if len(splitMsg) >= 2 {
			p.PrintAccount(peer, splitMsg[1])
			return true
		}
	default:
		fmt.Println("What you tried to send did not make sense")
		return false
	}
	return false
}
