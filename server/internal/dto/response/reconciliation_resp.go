package response

import "time"

type ReconciliationResp struct {
	ID              uint64    `json:"id"`
	AccountID       uint64    `json:"account_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	StartingBalance string    `json:"starting_balance"`
	EndingBalance   string    `json:"ending_balance"`
	BookBalance     string    `json:"book_balance"`
	Difference      string    `json:"difference"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ReconciliationEntryResp struct {
	ID               uint64 `json:"id"`
	ReconciliationID uint64 `json:"reconciliation_id"`
	TransactionID    uint64 `json:"transaction_id"`
	Matched          bool   `json:"matched"`
	AmountDifference string `json:"amount_difference"`
}