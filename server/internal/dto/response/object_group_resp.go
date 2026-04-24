package response

import "time"

type ObjectGroupResp struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	GroupableType string    `json:"groupable_type"`
	GroupableID   uint64    `json:"groupable_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}