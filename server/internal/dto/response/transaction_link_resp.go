package response

import "time"

type TransactionLinkResp struct {
	ID            uint64       `json:"id"`
	LinkTypeID    uint64       `json:"link_type_id"`
	SourceID      uint64       `json:"source_id"`
	DestinationID uint64       `json:"destination_id"`
	Comment       string       `json:"comment"`
	LinkType      LinkTypeResp `json:"link_type,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}
