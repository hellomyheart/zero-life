package response

import "time"

type LoginResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type ProfileResp struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Language  string    `json:"language"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}
