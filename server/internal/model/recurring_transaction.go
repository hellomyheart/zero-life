// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RecurrenceType 循环类型枚举
type RecurrenceType string

const (
	RecurrenceTypeDaily   RecurrenceType = "daily"   // 每日
	RecurrenceTypeWeekly  RecurrenceType = "weekly"  // 每周
	RecurrenceTypeMonthly RecurrenceType = "monthly" // 每月
	RecurrenceTypeYearly  RecurrenceType = "yearly"  // 每年
)

// RecurringTransaction 循环交易模型，对应recurring_transactions表
// 定义周期性交易模板，到期自动创建交易
type RecurringTransaction struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"not null;index" json:"user_id"`                   // 所属用户ID
	Description    string          `gorm:"not null;size:500" json:"description"`            // 交易描述
	Amount         decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`       // 交易金额
	SourceID       uint64          `gorm:"not null" json:"source_id"`                       // 源账户ID
	DestinationID  *uint64         `json:"destination_id"`                                  // 目标账户ID
	CategoryID     *uint64         `json:"category_id"`                                     // 分类ID
	Notes          string          `gorm:"type:text" json:"notes"`                          // 备注信息
	RecurrenceType RecurrenceType `gorm:"not null;size:20" json:"recurrence_type"`         // 循环类型
	RepeatEvery    int             `gorm:"default:1" json:"repeat_every"`                   // 循环间隔
	StartDate      time.Time       `gorm:"not null" json:"start_date"`                      // 开始日期
	EndDate        *time.Time      `json:"end_date,omitempty"`                              // 结束日期
	NextOccurrence time.Time       `gorm:"not null" json:"next_occurrence"`                 // 下次执行日期
	IsActive       bool            `gorm:"default:true" json:"is_active"`                   // 是否激活
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`                                  // 软删除时间
}

// TableName 指定表名
func (RecurringTransaction) TableName() string { return "recurring_transactions" }

// RecurringTransactionLog 循环交易日志模型，对应recurring_transaction_logs表
// 记录每次循环交易生成的交易
type RecurringTransactionLog struct {
	ID                  uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RecurringTransactionID uint64 `gorm:"not null;index" json:"recurring_transaction_id"`  // 循环交易ID
	TransactionID       uint64    `gorm:"not null" json:"transaction_id"`                  // 生成的交易ID
	OccurrenceDate      time.Time `gorm:"not null" json:"occurrence_date"`                 // 执行日期
	CreatedAt           time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (RecurringTransactionLog) TableName() string { return "recurring_transaction_logs" }
