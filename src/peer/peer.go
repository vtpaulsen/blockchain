package peer

import (
	"encoding/gob"
	"net"
	"strings"

	d "github.com/vtpaulsen/blockchain/v2/dataStructs"
	c "github.com/vtpaulsen/blockchain/v2/peer/connection"
	l "github.com/vtpaulsen/blockchain/v2/peer/ledger"
)

type Peer struct {
	ip          string
	port        int
	peers       map[int]string
	connections c.Connection
	ledger      l.Ledger
}

func SetIP(peer *Peer, ip string) {
	peer.ip = ip
}

func GetIP(peer *Peer) string {
	return peer.ip
}

func SetPort(peer *Peer, port int) {
	peer.port = port
}

func GetPort(peer *Peer) int {
	return peer.port
}

func AddToPeers(peer *Peer, port int, ip string) {
	peer.peers[port] = ip
}

func GetPeers(peer *Peer) map[int]string {
	return peer.peers
}

func GetConnections(peer *Peer) map[net.Conn]c.EncoderTuple {
	return c.GetConnections(&peer.connections)
}

func initializeConnections(peer *Peer) {
	c.InitializeConnections(&peer.connections)
}

func GetEncoder(peer *Peer, conn net.Conn) *gob.Encoder {
	return c.GetConnections(&peer.connections)[conn].Encoder
}

func GetDecoder(peer *Peer, conn net.Conn) *gob.Decoder {
	return c.GetConnections(&peer.connections)[conn].Decoder
}

func initializePeers(peer *Peer) {
	peer.peers = make(map[int]string)
}

func initializeLedger(peer *Peer) {

	l.InitializeLedger(&peer.ledger)

	// Test for signed transactions:
	// signedTransaction 101412213185284524806017175668973912724748971769560914346270211719807390186733 3 91805717802794537360572677055739513152082310523329363801232069588659394770961 3 100 67608142123523016537344783779315941816074659767719449455119487155241140580459

	// The exponent is 3 in all cases
	// signedTransaction n_from e_from n_to e_to amount d_from
	//                           (n, 3) : public                                            						             (n, d) : private
	l.AddToLedgers(&peer.ledger, "101412213185284524806017175668973912724748971769560914346270211719807390186733") // Secret key: n, 67608142123523016537344783779315941816074659767719449455119487155241140580459
	l.AddToLedgers(&peer.ledger, "91805717802794537360572677055739513152082310523329363801232069588659394770961")  // Secret key: n, 61203811868529691573715118037159675434316868581757632014150885561018447839555
	l.AddToLedgers(&peer.ledger, "216129593226809274476873310663022642783")                                        // Secret key: n, 144086395484539516298302447576604605147
	l.AddToLedgers(&peer.ledger, "302364108818827934793559482928635452821")                                        // Secret key: n, 201576072545885289839187010394497638067

	// Add 1000 to the first account
	l.GetLedger(&peer.ledger)["101412213185284524806017175668973912724748971769560914346270211719807390186733"] = 1000
}

func PrintAccount(peer *Peer, person string) {
	l.PrintAccount(&peer.ledger, person)
}

func PrintLedgers(peer *Peer) {
	l.PrintLedgers(&peer.ledger)
}

func AddToConnections(peer *Peer, conn net.Conn) {
	c.AddToConnections(&peer.connections, conn)
}

func AddToLedgers(peer *Peer, person string) {
	l.AddToLedgers(&peer.ledger, person)
}

func GetLedger(peer *Peer) map[string]int {
	return l.GetLedger(&peer.ledger)
}

func Transaction(peer *Peer, from string, to string, amount int) {
	l.Transaction(&peer.ledger, from, to, amount)
}

func SignedTransaction(peer *Peer, from string, to string, amount int) {
	splitFrom := strings.Split(from, ",")
	splitTo := strings.Split(to, ",")

	l.Transaction(&peer.ledger, splitFrom[0], splitTo[0], amount)
}

func Withdraw(peer *Peer, person string, amount int) {
	l.Withdraw(&peer.ledger, person, amount)
}

func Deposit(peer *Peer, person string, amount int) {
	l.Deposit(&peer.ledger, person, amount)
}

func IsWithdrawLegal(peer *Peer, person string, amount int) bool {
	//peer.ledger.mutex.Lock()
	if DoesPersonExist(peer, person) && amount > 0 {
		return GetLedger(peer)[person] >= amount
	}
	//defer peer.ledger.mutex.Unlock()
	return false
}

func DoesPersonExist(peer *Peer, person string) bool {
	//peer.ledger.mutex.Lock()
	if _, ok := GetLedger(peer)[person]; ok {
		//defer peer.ledger.mutex.Unlock()
		return true
	}
	return false
}

func InitializePeer() Peer {

	gob.Register(d.Message{})
	gob.Register(d.TransactionStruct{})
	gob.Register(d.SignedTransactionStruct{})

	var peer Peer
	SetIP(&peer, "")
	SetPort(&peer, 0)
	initializePeers(&peer)
	initializeConnections(&peer)
	initializeLedger(&peer)

	return peer
}
