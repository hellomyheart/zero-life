// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Currency 货币模型，对应currencies表
// 管理系统支持的货币类型和汇率信息
type Currency struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Code          string `gorm:"uniqueIndex;not null;size:3" json:"code"`       // 货币代码，如CNY、USD
	Name          string `gorm:"not null;size:100" json:"name"`                  // 货币名称
	Symbol        string `gorm:"not null;size:10" json:"symbol"`                 // 货币符号，如¥、$
	DecimalPlaces int    `gorm:"default:2" json:"decimal_places"`                // 小数位数，如JPY为0
	IsEnabled     bool   `gorm:"default:true" json:"is_enabled"`                 // 是否启用
	IsDefault     bool   `gorm:"default:false" json:"is_default"`                // 是否为默认货币
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (Currency) TableName() string { return "currencies" }

// ExchangeRate 汇率模型，对应exchange_rates表
// 记录不同货币之间的兑换比率
type ExchangeRate struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromCurrencyID uint64          `gorm:"not null;index:idx_currency_pair,unique" json:"from_currency_id"` // 源货币ID
	ToCurrencyID   uint64          `gorm:"not null;index:idx_currency_pair,unique" json:"to_currency_id"`   // 目标货币ID
	Rate           decimal.Decimal `gorm:"type:decimal(19,8);not null" json:"rate"`                         // 兑换比率
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (ExchangeRate) TableName() string { return "exchange_rates" }
