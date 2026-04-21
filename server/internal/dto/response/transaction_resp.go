package response

import "time"

type TransactionResp struct {
	ID            uint64        `json:"id"`
	Type          string        `json:"type"`
	Date          time.Time     `json:"date"`
	Description   string        `json:"description"`
	Amount        string        `json:"amount"`
	SourceID      uint64        `json:"source_id"`
	Source        AccountResp   `json:"source"`
	DestinationID *uint64       `json:"destination_id"`
	Destination   *AccountResp  `json:"destination,omitempty"`
	CategoryID    *uint64       `json:"category_id"`
	Category      *CategoryResp `json:"category,omitempty"`
	Notes         string        `json:"notes"`
	Tags          []TagResp     `json:"tags"`
	Splits        []SplitResp   `json:"splits,omitempty"`
	BillID        *uint64       `json:"bill_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type SplitResp struct {
	ID         uint64        `json:"id"`
	Amount     string        `json:"amount"`
	CategoryID *uint64       `json:"category_id"`
	Category   *CategoryResp `json:"category,omitempty"`
	Tags       []TagResp     `json:"tags"`
	Notes      string        `json:"notes"`
}
