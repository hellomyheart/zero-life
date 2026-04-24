package response

import "time"

// UserResp 用户响应
type UserResp struct {
	ID         uint64    `json:"id"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	Language   string    `json:"language"`
	Timezone   string    `json:"timezone"`
	Role       string    `json:"role"`
	MFAEnabled bool      `json:"mfa_enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
