// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// WebhookResp Webhook响应
type WebhookResp struct {
	ID              uint64     `json:"id"`                          // Webhook ID
	Name            string     `json:"name"`                        // Webhook名称
	URL             string     `json:"url"`                         // Webhook接收URL
	Trigger         string     `json:"trigger"`                     // 触发事件类型
	IsActive        bool       `json:"is_active"`                   // 是否启用
	LastDeliveredAt *time.Time `json:"last_delivered_at,omitempty"` // 最后投递时间
	CreatedAt       time.Time  `json:"created_at"`                  // 创建时间
	UpdatedAt       time.Time  `json:"updated_at"`                  // 更新时间
}

// WebhookDeliveryResp Webhook投递记录响应
type WebhookDeliveryResp struct {
	ID           uint64     `json:"id"`                         // 投递记录ID
	WebhookID    uint64     `json:"webhook_id"`                 // Webhook ID
	StatusCode   *int       `json:"status_code,omitempty"`      // HTTP响应状态码
	ErrorMessage *string    `json:"error_message,omitempty"`    // 错误信息
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`     // 投递时间
	CreatedAt    time.Time  `json:"created_at"`                 // 创建时间
}