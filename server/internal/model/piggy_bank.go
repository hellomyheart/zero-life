package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PiggyBank struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64          `gorm:"not null;index" json:"user_id"`
	Name          string          `gorm:"not null;size:100" json:"name"`
	TargetAmount  decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"target_amount"`
	CurrentAmount decimal.Decimal `gorm:"type:decimal(19,4);not null;default:0" json:"current_amount"`
	AccountID     uint64          `gorm:"not null" json:"account_id"`
	Order         int             `gorm:"default:0" json:"order"`
	TargetDate    *time.Time      `json:"target_date,omitempty"`
	Notes         string          `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`

	Account Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Events  []PiggyEvent `gorm:"foreignKey:PiggyBankID" json:"events,omitempty"`
}

func (PiggyBank) TableName() string { return "piggy_banks" }

type PiggyEvent struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PiggyBankID   uint64          `gorm:"not null;index" json:"piggy_bank_id"`
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	TransactionID *uint64         `json:"transaction_id,omitempty"`
	Note          string          `gorm:"type:text" json:"note"`
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`
}

func (PiggyEvent) TableName() string { return "piggy_events" }
