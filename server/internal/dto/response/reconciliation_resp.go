package response

import "time"

type ReconciliationDataResp struct {
	AccountID    uint64            `json:"account_id"`
	StartDate    time.Time         `json:"start_date"`
	EndDate      time.Time         `json:"end_date"`
	StartBalance string            `json:"start_balance"`
	EndBalance   string            `json:"end_balance"`
	Transactions []TransactionResp `json:"transactions"`
}

type ReconciliationResp struct {
	ID               uint64    `json:"id"`
	AccountID        uint64    `json:"account_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	StartBalance     string    `json:"start_balance"`
	EndBalance       string    `json:"end_balance"`
	SubmittedBalance string    `json:"submitted_balance"`
	Difference       string    `json:"difference"`
	CreatedAt        time.Time `json:"created_at"`
}
