package response

import "time"

type WebhookResp struct {
	ID              uint64     `json:"id"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	Trigger         string     `json:"trigger"`
	IsActive        bool       `json:"is_active"`
	LastDeliveredAt *time.Time `json:"last_delivered_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type WebhookDeliveryResp struct {
	ID           uint64     `json:"id"`
	WebhookID    uint64     `json:"webhook_id"`
	StatusCode   *int       `json:"status_code,omitempty"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}