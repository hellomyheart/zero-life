package model

import (
	"time"

	"gorm.io/gorm"
)

type WebhookTrigger string

const (
	WebhookTriggerTransactionCreate WebhookTrigger = "TRANSACTION_CREATE"
	WebhookTriggerTransactionUpdate WebhookTrigger = "TRANSACTION_UPDATE"
	WebhookTriggerTransactionDelete WebhookTrigger = "TRANSACTION_DELETE"
)

type WebhookMessageStatus string

const (
	WebhookMessagePending WebhookMessageStatus = "pending"
	WebhookMessageSuccess WebhookMessageStatus = "success"
	WebhookMessageFailed  WebhookMessageStatus = "failed"
)

type Webhook struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	URL       string         `gorm:"not null;size:500" json:"url"`
	Trigger   WebhookTrigger `gorm:"not null;size:50" json:"trigger"`
	Secret    string         `gorm:"size:255" json:"secret"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Webhook) TableName() string { return "webhooks" }

type WebhookMessage struct {
	ID           uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	WebhookID    uint64              `gorm:"not null;index" json:"webhook_id"`
	RequestBody  string             `gorm:"type:text" json:"request_body"`
	ResponseCode *int               `json:"response_code"`
	ResponseBody string             `gorm:"type:text" json:"response_body"`
	Attempts     int                `gorm:"default:0" json:"attempts"`
	Status       WebhookMessageStatus `gorm:"size:20;default:pending" json:"status"`
	SentAt       *time.Time         `json:"sent_at"`
	CreatedAt    time.Time          `gorm:"not null" json:"created_at"`
}

func (WebhookMessage) TableName() string { return "webhook_messages" }
