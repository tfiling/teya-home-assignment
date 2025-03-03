package ledger

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction represents an internal model for transaction entity
type Transaction struct {
	ID         uint64          `json:"id"`
	Amount     decimal.Decimal `json:"amount"`
	ExternalID uuid.UUID       `json:"external_id"`
}
