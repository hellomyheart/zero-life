// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccountType 账户类型枚举
type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"     // 资产账户：银行账户、现金等
	AccountTypeExpense   AccountType = "expense"   // 支出账户：用于记录支出
	AccountTypeRevenue   AccountType = "revenue"   // 收入账户：用于记录收入
	AccountTypeLiability AccountType = "liability" // 负债账户：信用卡、贷款等
)

// Account 账户模型，对应accounts表
// 管理用户的各类财务账户，包括资产、负债、收入、支出账户
type Account struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`                  // 账户ID，主键自增
	UserID         uint64          `gorm:"not null;index" json:"user_id"`                       // 所属用户ID
	Name           string          `gorm:"not null;size:255" json:"name"`                       // 账户名称
	Type           AccountType     `gorm:"not null;size:20" json:"type"`                        // 账户类型（asset/expense/revenue/liability）
	CurrencyID     uint64          `gorm:"not null" json:"currency_id"`                         // 货币ID
	InitialBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"initial_balance"` // 初始余额
	CurrentBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"current_balance"` // 当前余额（自动计算）
	IsVirtual      bool            `gorm:"default:false" json:"is_virtual"`                     // 是否为虚拟账户
	Notes          string          `gorm:"type:text" json:"notes"`                              // 备注信息
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`                          // 创建时间
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`                          // 更新时间
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`                                      // 软删除时间

	Currency Currency `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"` // 关联货币信息
}

// TableName 指定表名
func (Account) TableName() string { return "accounts" }
