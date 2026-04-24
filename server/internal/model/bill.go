// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RepeatRule 重复规则枚举
type RepeatRule string

const (
	RepeatRuleDaily   RepeatRule = "daily"   // 每日
	RepeatRuleWeekly  RepeatRule = "weekly"  // 每周
	RepeatRuleMonthly RepeatRule = "monthly" // 每月
	RepeatRuleYearly  RepeatRule = "yearly"  // 每年
)

// Bill 账单模型，对应bills表
// 管理周期性账单，自动提醒到期付款
type Bill struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64          `gorm:"not null;index" json:"user_id"`                // 所属用户ID
	Name       string          `gorm:"not null;size:100" json:"name"`                // 账单名称
	Amount     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`    // 账单金额
	RepeatRule RepeatRule      `gorm:"not null;size:20" json:"repeat_rule"`          // 重复规则（daily/weekly/monthly/yearly）
	NextDue    time.Time       `gorm:"not null" json:"next_due"`                     // 下次到期日期
	SourceID   *uint64         `json:"source_id"`                                    // 支出账户ID
	CategoryID *uint64         `json:"category_id"`                                  // 分类ID
	Notes      string          `gorm:"type:text" json:"notes"`                       // 备注信息
	CreatedAt  time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`                               // 软删除时间
}

// TableName 指定表名
func (Bill) TableName() string { return "bills" }
