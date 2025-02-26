package ledger

import (
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	TransactionHistory []Transaction
	transactionIdSeq   atomic.Uint64
}

func NewLedger() (*Ledger, error) {
	l := &Ledger{
		TransactionHistory: make([]Transaction, 0),
	}
	l.transactionIdSeq.Store(0)
	return l, nil
}

func (l *Ledger) AddTransaction(amount decimal.Decimal) {
	newTransaction := Transaction{
		ID:         l.getNewID(),
		Amount:     amount,
		ExternalID: uuid.New(),
	}
	l.TransactionHistory = append(l.TransactionHistory, newTransaction)
}

func (l *Ledger) GetBalance() decimal.Decimal {
	balance := decimal.NewFromFloat(0)
	for _, transaction := range l.TransactionHistory {
		balance = balance.Add(transaction.Amount)
	}
	return balance
}
func (l *Ledger) getNewID() uint64 {
	return l.transactionIdSeq.Add(1)
}
