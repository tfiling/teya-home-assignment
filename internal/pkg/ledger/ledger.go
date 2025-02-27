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

func (l *Ledger) AddTransaction(amount decimal.Decimal) error {
	newTransaction := Transaction{
		ID:         l.getNewID(),
		Amount:     amount,
		ExternalID: uuid.New(),
	}
	l.TransactionHistory = append(l.TransactionHistory, newTransaction)
	return nil
}

func (l *Ledger) GetBalance() (decimal.Decimal, error) {
	//TODO - optimize performance for large ledger
	balance := decimal.NewFromFloat(0)
	for _, transaction := range l.TransactionHistory {
		balance = balance.Add(transaction.Amount)
	}
	return balance, nil
}

func (l *Ledger) GetTransactionHistory(offset, limit int) ([]Transaction, error) {
	if offset > len(l.TransactionHistory) {
		//TODO - consider returning an error here
		return []Transaction{}, nil
	}

	endIndex := offset + limit
	if endIndex > len(l.TransactionHistory) {
		endIndex = len(l.TransactionHistory)
	}

	return l.TransactionHistory[offset:endIndex], nil
}

func (l *Ledger) getNewID() uint64 {
	return l.transactionIdSeq.Add(1)
}
