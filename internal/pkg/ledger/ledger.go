package ledger

import (
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	TransactionHistory   []Transaction
	transactionIdSeq     atomic.Uint64
	cachedBalance        decimal.Decimal
	cachedBalanceTillIdx int64
}

func NewLedger() (*Ledger, error) {
	l := &Ledger{
		TransactionHistory: make([]Transaction, 0),
		//cachedBalanceTillIdx is inclusive and at this point we did not cache the first transaction
		cachedBalanceTillIdx: -1,
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
	balance := l.cachedBalance
	for i := l.cachedBalanceTillIdx + 1; i < int64(len(l.TransactionHistory)); i++ {
		balance = balance.Add(l.TransactionHistory[i].Amount)
		//Cache balance that was already calculated
		l.cachedBalanceTillIdx = i
		l.cachedBalance = balance
	}
	return balance, nil
}

func (l *Ledger) GetTransactionHistory(offset, limit int) ([]Transaction, error) {
	if offset > len(l.TransactionHistory) {
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
