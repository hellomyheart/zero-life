package response

type PiggyBankResp struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`
	TargetAmount  string  `json:"target_amount"`
	CurrentAmount string  `json:"current_amount"`
	AccountID     uint64  `json:"account_id"`
	TargetDate    *string `json:"target_date,omitempty"`
	Notes         string  `json:"notes"`
	Percentage    float64 `json:"percentage"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type PiggyEventResp struct {
	ID            uint64  `json:"id"`
	PiggyBankID   uint64  `json:"piggy_bank_id"`
	Amount        string  `json:"amount"`
	TransactionID *uint64 `json:"transaction_id,omitempty"`
	Note          string  `json:"note"`
	CreatedAt     string  `json:"created_at"`
}
