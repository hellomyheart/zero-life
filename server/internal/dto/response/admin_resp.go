package response

import "time"

type AdminUserResp struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Role      string    `json:"role"`
	Language  string    `json:"language"`
	Timezone  string    `json:"timezone"`
	MFAEnabled bool     `json:"mfa_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminConfigurationResp struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminTestEmailResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
