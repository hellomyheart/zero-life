package response

import "time"

type TagResp struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Color           string    `json:"color"`
	TransactionCount int64    `json:"transaction_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
