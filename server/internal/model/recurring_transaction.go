package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RecurrenceType string

const (
	RecurrenceTypeDaily   RecurrenceType = "daily"
	RecurrenceTypeWeekly  RecurrenceType = "weekly"
	RecurrenceTypeMonthly RecurrenceType = "monthly"
	RecurrenceTypeYearly  RecurrenceType = "yearly"
)

type RecurringTransaction struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"not null;index" json:"user_id"`
	Description    string          `gorm:"not null;size:500" json:"description"`
	Amount         decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	SourceID       uint64          `gorm:"not null" json:"source_id"`
	DestinationID  *uint64         `json:"destination_id"`
	CategoryID     *uint64         `json:"category_id"`
	Notes          string          `gorm:"type:text" json:"notes"`
	RecurrenceType RecurrenceType `gorm:"not null;size:20" json:"recurrence_type"`
	RepeatEvery    int             `gorm:"default:1" json:"repeat_every"`
	StartDate      time.Time       `gorm:"not null" json:"start_date"`
	EndDate        *time.Time      `json:"end_date,omitempty"`
	NextOccurrence time.Time       `gorm:"not null" json:"next_occurrence"`
	IsActive       bool            `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (RecurringTransaction) TableName() string { return "recurring_transactions" }

type RecurringTransactionLog struct {
	ID                  uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RecurringTransactionID uint64 `gorm:"not null;index" json:"recurring_transaction_id"`
	TransactionID       uint64    `gorm:"not null" json:"transaction_id"`
	OccurrenceDate      time.Time `gorm:"not null" json:"occurrence_date"`
	CreatedAt           time.Time `gorm:"not null" json:"created_at"`
}

func (RecurringTransactionLog) TableName() string { return "recurring_transaction_logs" }