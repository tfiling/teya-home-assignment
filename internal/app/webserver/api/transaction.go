package api

import (
	"teya_home_assignment/internal/pkg/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID     uuid.UUID       `json:"id"`
	Amount decimal.Decimal `json:"amount"`
}

type NewTransactionReqBody struct {
	Amount string `json:"amount"`
}

type GetBalanceRespBody struct {
	Balance string `json:"balance"`
}

func FromTransactionModel(transaction ledger.Transaction) Transaction {
	return Transaction{
		ID:     transaction.ExternalID,
		Amount: transaction.Amount,
	}
}
