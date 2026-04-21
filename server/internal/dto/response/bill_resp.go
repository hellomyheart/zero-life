package response

import "time"

type BillResp struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	Amount     string    `json:"amount"`
	RepeatRule string    `json:"repeat_rule"`
	NextDue    time.Time `json:"next_due"`
	SourceID   *uint64   `json:"source_id"`
	CategoryID *uint64   `json:"category_id"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
