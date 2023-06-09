package dataStructs

type Message struct {
	TransactionType string
	Payload         any
}

type TransactionStruct struct {
	ID     string
	From   string
	To     string
	Amount int
}

type SignedTransactionStruct struct {
	ID        string
	From      string
	To        string
	Amount    int
	Signature string
}

const (
	TRANSACTION        = "transaction"
	SIGNED_TRANSACTION = "signedTransaction"
	JOIN_MSG           = "join"
	ADD_MSG            = "add"
	REQUEST_MSG        = "request"
	REQUEST_REPLY      = "requestreply"
	ADD_TO_LEDGER      = "addPersonToLedger"
	WITHDRAW           = "withdraw"
	DEPOSIT            = "deposit"
	PRINTLEDGER        = "printledger"
	MSG                = "msg"
	GET_ACCOUNT        = "getAccount"
)
