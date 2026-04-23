package response

import "time"

type LinkTypeResp struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Outward       string    `json:"outward"`
	Inward        string    `json:"inward"`
	IsDirectional bool      `json:"is_directional"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
