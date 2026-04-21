package response

import "time"

type AccountResp struct {
	ID             uint64       `json:"id"`
	Name           string       `json:"name"`
	Type           string       `json:"type"`
	CurrencyID     uint64       `json:"currency_id"`
	Currency       CurrencyResp `json:"currency"`
	InitialBalance string       `json:"initial_balance"`
	CurrentBalance string       `json:"current_balance"`
	IsVirtual      bool         `json:"is_virtual"`
	Notes          string       `json:"notes"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
