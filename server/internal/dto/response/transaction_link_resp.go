package response

import "time"

type TransactionLinkResp struct {
	ID              uint64    `json:"id"`
	TransactionID   uint64    `json:"transaction_id"`
	LinkType        string    `json:"link_type"`
	LinkedJournalID uint64    `json:"linked_journal_id"`
	CreatedAt       time.Time `json:"created_at"`
}