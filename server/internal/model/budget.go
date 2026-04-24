// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BudgetPeriod 预算周期类型
type BudgetPeriod string

const (
	BudgetPeriodMonthly BudgetPeriod = "monthly" // 月度预算
	BudgetPeriodYearly  BudgetPeriod = "yearly"  // 年度预算
)

// Budget 预算模型，对应budgets表
// 设定分类支出上限，跟踪预算使用情况
type Budget struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64          `gorm:"not null;index" json:"user_id"`           // 所属用户ID
	Name      string          `gorm:"not null;size:100" json:"name"`           // 预算名称
	Amount    decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"` // 预算金额上限
	Period    BudgetPeriod    `gorm:"not null;size:20" json:"period"`          // 预算周期（monthly/yearly）
	IsEnabled bool            `gorm:"default:true" json:"is_enabled"`          // 是否启用
	CreatedAt time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`                          // 软删除时间

	Categories []Category `gorm:"many2many:budget_categories" json:"categories,omitempty"` // 关联分类列表
}

// TableName 指定表名
func (Budget) TableName() string { return "budgets" }

// BudgetCategory 预算-分类关联模型，对应budget_categories表
type BudgetCategory struct {
	BudgetID   uint64 `gorm:"primaryKey" json:"budget_id"`   // 预算ID
	CategoryID uint64 `gorm:"primaryKey" json:"category_id"` // 分类ID
}

// TableName 指定表名
func (BudgetCategory) TableName() string { return "budget_categories" }

// BudgetHistory 预算历史模型，对应budget_history表
// 记录每个预算周期的实际支出和预算金额
type BudgetHistory struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BudgetID    uint64          `gorm:"not null;index" json:"budget_id"`           // 预算ID
	PeriodStart time.Time       `gorm:"not null" json:"period_start"`              // 周期开始日期
	PeriodEnd   time.Time       `gorm:"not null" json:"period_end"`                // 周期结束日期
	Amount      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"` // 预算金额
	Spent       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"spent"`  // 实际支出金额
	CreatedAt   time.Time       `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (BudgetHistory) TableName() string { return "budget_history" }
