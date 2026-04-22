package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RepeatRule string

const (
	RepeatRuleDaily   RepeatRule = "daily"
	RepeatRuleWeekly  RepeatRule = "weekly"
	RepeatRuleMonthly RepeatRule = "monthly"
	RepeatRuleYearly  RepeatRule = "yearly"
)

type Bill struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64          `gorm:"not null;index" json:"user_id"`
	Name       string          `gorm:"not null;size:100" json:"name"`
	Amount     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	RepeatRule RepeatRule      `gorm:"not null;size:20" json:"repeat_rule"`
	NextDue    time.Time       `gorm:"not null" json:"next_due"`
	SourceID   *uint64         `json:"source_id"`
	CategoryID *uint64         `json:"category_id"`
	Notes      string          `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Bill) TableName() string { return "bills" }
