// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// WebhookTrigger Webhook触发事件类型
type WebhookTrigger string

const (
	WebhookTriggerTransactionCreate WebhookTrigger = "TRANSACTION_CREATE" // 交易创建
	WebhookTriggerTransactionUpdate WebhookTrigger = "TRANSACTION_UPDATE" // 交易更新
	WebhookTriggerTransactionDelete WebhookTrigger = "TRANSACTION_DELETE" // 交易删除
)

// WebhookMessageStatus Webhook消息状态
type WebhookMessageStatus string

const (
	WebhookMessagePending WebhookMessageStatus = "pending" // 待发送
	WebhookMessageSuccess WebhookMessageStatus = "success" // 发送成功
	WebhookMessageFailed  WebhookMessageStatus = "failed"  // 发送失败
)

// Webhook Webhook模型，对应webhooks表
// 配置外部回调地址，在特定事件发生时推送通知
type Webhook struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint64         `gorm:"not null;index" json:"user_id"`                // 所属用户ID
	Name            string         `gorm:"not null;size:100" json:"name"`                // Webhook名称
	URL             string         `gorm:"not null;size:500" json:"url"`                 // 回调URL
	Trigger         WebhookTrigger `gorm:"not null;size:50" json:"trigger"`              // 触发事件类型
	Secret          string         `gorm:"size:255" json:"secret"`                       // 签名密钥
	IsActive        bool           `gorm:"default:true" json:"is_active"`                // 是否启用
	LastDeliveredAt *time.Time     `json:"last_delivered_at,omitempty"`                  // 最后投递时间
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除时间
}

// TableName 指定表名
func (Webhook) TableName() string { return "webhooks" }

// WebhookMessage Webhook消息模型，对应webhook_messages表
// 记录每次Webhook推送的请求和响应
type WebhookMessage struct {
	ID           uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	WebhookID    uint64              `gorm:"not null;index" json:"webhook_id"`              // Webhook ID
	RequestBody  string             `gorm:"type:text" json:"request_body"`                 // 请求体
	ResponseCode *int               `json:"response_code"`                                 // HTTP响应状态码
	ResponseBody string             `gorm:"type:text" json:"response_body"`                // 响应体
	Attempts     int                `gorm:"default:0" json:"attempts"`                     // 重试次数
	Status       WebhookMessageStatus `gorm:"size:20;default:pending" json:"status"`         // 消息状态
	SentAt       *time.Time         `json:"sent_at"`                                       // 发送时间
	CreatedAt    time.Time          `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (WebhookMessage) TableName() string { return "webhook_messages" }

// WebhookDelivery Webhook投递记录模型，对应webhook_deliveries表
type WebhookDelivery struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WebhookID    uint64     `gorm:"not null;index" json:"webhook_id"`     // Webhook ID
	Payload      string     `gorm:"type:text" json:"payload"`             // 推送内容
	StatusCode   *int       `json:"status_code,omitempty"`                // HTTP状态码
	ErrorMessage *string    `json:"error_message,omitempty"`              // 错误信息
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`               // 投递时间
	CreatedAt    time.Time  `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (WebhookDelivery) TableName() string { return "webhook_deliveries" }
