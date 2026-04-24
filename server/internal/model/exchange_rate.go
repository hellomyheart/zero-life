// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ExchangeRate 汇率模型，对应currency_exchange_rates表
// 记录不同货币之间的汇率信息，支持按日期查询历史汇率
type ExchangeRate struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`              // 汇率ID，主键自增
	UserID         uint64          `gorm:"not null;index" json:"user_id"`                   // 所属用户ID
	UserGroupID    uint64          `gorm:"index" json:"user_group_id"`                      // 用户组ID（可选）
	FromCurrencyID uint64          `gorm:"not null;index" json:"from_currency_id"`          // 源货币ID
	ToCurrencyID   uint64          `gorm:"not null;index" json:"to_currency_id"`            // 目标货币ID
	Date           time.Time       `gorm:"not null;type:date;index" json:"date"`            // 汇率日期
	Rate           decimal.Decimal `gorm:"not null;type:decimal(19,12)" json:"rate"`        // 汇率值
	UserRate       decimal.Decimal `gorm:"type:decimal(19,12)" json:"user_rate,omitempty"`  // 用户自定义汇率（可选）
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`                      // 创建时间
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`                      // 更新时间
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`                                  // 软删除时间

	FromCurrency Currency `gorm:"foreignKey:FromCurrencyID" json:"from_currency,omitempty"` // 源货币信息
	ToCurrency   Currency `gorm:"foreignKey:ToCurrencyID" json:"to_currency,omitempty"`     // 目标货币信息
}

// TableName 指定表名
func (ExchangeRate) TableName() string { return "currency_exchange_rates" }
