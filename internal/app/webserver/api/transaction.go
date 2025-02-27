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
	Amount string `json:"amount" validate:"required,number"`
}

type GetBalanceRespBody struct {
	Balance string `json:"balance"`
}

type PaginatedTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Pagination   Pagination    `json:"pagination"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func FromTransactionModel(transaction ledger.Transaction) Transaction {
	return Transaction{
		ID:     transaction.ExternalID,
		Amount: transaction.Amount,
	}
}
