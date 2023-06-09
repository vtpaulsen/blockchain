package ledger

import (
	"fmt"
	"sync"
)

type Ledger struct {
	ledger map[string]int
	mutex  sync.Mutex
}

func lockLedgerMutex(ledger *Ledger) {
	ledger.mutex.Lock()
}

func unlockLedgerMutex(ledger *Ledger) {
	ledger.mutex.Unlock()
}

func GetLedger(ledger *Ledger) map[string]int {
	return ledger.ledger
}

func setLedger(ledger *Ledger, person string, amount int) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	ledger.ledger[person] = amount
}

func InitializeLedger(ledger *Ledger) {
	ledger.ledger = make(map[string]int)
}

func AddToLedgers(ledger *Ledger, person string) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	ledger.ledger[person] = 0
}

func PrintLedgers(ledger *Ledger) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	for key, value := range ledger.ledger {
		fmt.Println(key, value)
	}
}

func PrintAccount(ledger *Ledger, person string) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	fmt.Println(person, GetLedger(ledger)[person])
}

func Transaction(ledger *Ledger, from string, to string, amount int) {
	//lockLedgerMutex(ledger)
	//defer unlockLedgerMutex(ledger)
	setLedger(ledger, from, GetLedger(ledger)[from]-amount)
	setLedger(ledger, to, GetLedger(ledger)[to]+amount)
}

func Withdraw(ledger *Ledger, person string, amount int) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	setLedger(ledger, person, GetLedger(ledger)[person]-amount)
}

func Deposit(ledger *Ledger, person string, amount int) {
	lockLedgerMutex(ledger)
	defer unlockLedgerMutex(ledger)
	setLedger(ledger, person, GetLedger(ledger)[person]+amount)
}
