package response

import "time"

type WebhookResp struct {
	ID        uint64    `json:"id"`
	URL       string    `json:"url"`
	Trigger   string    `json:"trigger"`
	Secret    string    `json:"secret"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WebhookMessageResp struct {
	ID           uint64    `json:"id"`
	WebhookID    uint64    `json:"webhook_id"`
	RequestBody  string    `json:"request_body"`
	ResponseCode *int      `json:"response_code"`
	ResponseBody string    `json:"response_body"`
	Attempts     int       `json:"attempts"`
	Status       string    `json:"status"`
	SentAt       *time.Time `json:"sent_at"`
	CreatedAt    time.Time `json:"created_at"`
}
