package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Currency struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Code          string         `gorm:"uniqueIndex;not null;size:3" json:"code"`
	Name          string         `gorm:"not null;size:100" json:"name"`
	Symbol        string         `gorm:"not null;size:10" json:"symbol"`
	DecimalPlaces int           `gorm:"default:2" json:"decimal_places"`
	IsEnabled     bool          `gorm:"default:true" json:"is_enabled"`
	IsDefault     bool          `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"not null" json:"updated_at"`
}

func (Currency) TableName() string { return "currencies" }

type ExchangeRate struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromCurrencyID uint64          `gorm:"not null;index:idx_currency_pair,unique" json:"from_currency_id"`
	ToCurrencyID   uint64          `gorm:"not null;index:idx_currency_pair,unique" json:"to_currency_id"`
	Rate           decimal.Decimal `gorm:"type:decimal(19,8);not null" json:"rate"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
}

func (ExchangeRate) TableName() string { return "exchange_rates" }
