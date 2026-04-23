package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RepeatFreq string

const (
	RepeatFreqDaily   RepeatFreq = "daily"
	RepeatFreqWeekly  RepeatFreq = "weekly"
	RepeatFreqMonthly RepeatFreq = "monthly"
	RepeatFreqYearly  RepeatFreq = "yearly"
)

type Recurrence struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"not null;index" json:"user_id"`
	Title          string          `gorm:"not null;size:255" json:"title"`
	Type           TransactionType `gorm:"not null;size:20" json:"type"`
	Amount         decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	SourceID       uint64          `gorm:"not null" json:"source_id"`
	DestinationID  *uint64         `json:"destination_id"`
	CategoryID     *uint64         `json:"category_id"`
	Description    string          `gorm:"size:500" json:"description"`
	Notes          string          `gorm:"type:text" json:"notes"`
	TagNames       string          `gorm:"size:500" json:"tag_names"`
	RepeatFreq     RepeatFreq      `gorm:"not null;size:20" json:"repeat_freq"`
	RepeatInterval int             `gorm:"not null;default:1" json:"repeat_interval"`
	NextDate       time.Time       `gorm:"not null" json:"next_date"`
	EndDate        *time.Time      `json:"end_date"`
	Repetitions    int             `gorm:"default:0" json:"repetitions"`
	MaxRepetitions *int            `json:"max_repetitions"`
	IsActive       bool            `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Recurrence) TableName() string { return "recurrences" }
