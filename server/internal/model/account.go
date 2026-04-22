package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeExpense   AccountType = "expense"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeLiability AccountType = "liability"
)

type Account struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"not null;index" json:"user_id"`
	Name           string          `gorm:"not null;size:255" json:"name"`
	Type           AccountType     `gorm:"not null;size:20" json:"type"`
	CurrencyID     uint64          `gorm:"not null" json:"currency_id"`
	InitialBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"initial_balance"`
	CurrentBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"current_balance"`
	IsVirtual      bool            `gorm:"default:false" json:"is_virtual"`
	Notes          string          `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`

	Currency Currency `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"`
}

func (Account) TableName() string { return "accounts" }
