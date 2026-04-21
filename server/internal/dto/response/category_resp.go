package response

import "time"

type CategoryResp struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	ParentID  *uint64        `json:"parent_id"`
	Icon      string         `json:"icon"`
	Notes     string         `json:"notes"`
	SortOrder int           `json:"sort_order"`
	Children  []CategoryResp `json:"children,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
