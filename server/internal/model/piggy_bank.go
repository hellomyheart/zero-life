// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PiggyBank 储蓄罐模型，对应piggy_banks表
// 设定储蓄目标，跟踪储蓄进度
type PiggyBank struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64          `gorm:"not null;index" json:"user_id"`                    // 所属用户ID
	Name          string          `gorm:"not null;size:100" json:"name"`                    // 储蓄罐名称
	TargetAmount  decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"target_amount"` // 目标金额
	CurrentAmount decimal.Decimal `gorm:"type:decimal(19,4);not null;default:0" json:"current_amount"` // 当前已存金额
	AccountID     uint64          `gorm:"not null" json:"account_id"`                       // 关联账户ID
	Order         int             `gorm:"default:0" json:"order"`                           // 排序序号
	TargetDate    *time.Time      `json:"target_date,omitempty"`                            // 目标达成日期
	Notes         string          `gorm:"type:text" json:"notes"`                           // 备注信息
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`                                   // 软删除时间

	Account Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"` // 关联账户
	Events  []PiggyEvent `gorm:"foreignKey:PiggyBankID" json:"events,omitempty"` // 存取记录列表
}

// TableName 指定表名
func (PiggyBank) TableName() string { return "piggy_banks" }

// PiggyEvent 储蓄事件模型，对应piggy_events表
// 记录储蓄罐的每次存入或取出操作
type PiggyEvent struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PiggyBankID   uint64          `gorm:"not null;index" json:"piggy_bank_id"`              // 储蓄罐ID
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`        // 金额（正数存入，负数取出）
	TransactionID *uint64         `json:"transaction_id,omitempty"`                         // 关联交易ID
	Note          string          `gorm:"type:text" json:"note"`                            // 备注信息
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (PiggyEvent) TableName() string { return "piggy_events" }
