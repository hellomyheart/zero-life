package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Reconciliation struct {
	ID               uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64          `gorm:"not null;index" json:"user_id"`
	AccountID        uint64          `gorm:"not null;index" json:"account_id"`
	StartDate        time.Time       `gorm:"not null" json:"start_date"`
	EndDate          time.Time       `gorm:"not null" json:"end_date"`
	StartBalance     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"start_balance"`
	EndBalance       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"end_balance"`
	SubmittedBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"submitted_balance"`
	Difference       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"difference"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
}

func (Reconciliation) TableName() string { return "reconciliations" }
