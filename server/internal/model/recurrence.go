// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RepeatFreq 重复频率枚举
type RepeatFreq string

const (
	RepeatFreqDaily   RepeatFreq = "daily"   // 每日
	RepeatFreqWeekly  RepeatFreq = "weekly"  // 每周
	RepeatFreqMonthly RepeatFreq = "monthly" // 每月
	RepeatFreqYearly  RepeatFreq = "yearly"  // 每年
)

// Recurrence 定期交易模型，对应recurrences表
// 定义周期性交易模板，到期自动创建交易
type Recurrence struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"not null;index" json:"user_id"`                   // 所属用户ID
	Title          string          `gorm:"not null;size:255" json:"title"`                  // 定期交易名称
	Type           TransactionType `gorm:"not null;size:20" json:"type"`                    // 交易类型
	Amount         decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`       // 交易金额
	SourceID       uint64          `gorm:"not null" json:"source_id"`                       // 源账户ID
	DestinationID  *uint64         `json:"destination_id"`                                  // 目标账户ID
	CategoryID     *uint64         `json:"category_id"`                                     // 分类ID
	Description    string          `gorm:"size:500" json:"description"`                     // 交易描述
	Notes          string          `gorm:"type:text" json:"notes"`                          // 备注信息
	TagNames       string          `gorm:"size:500" json:"tag_names"`                       // 标签名称（逗号分隔）
	RepeatFreq     RepeatFreq      `gorm:"not null;size:20" json:"repeat_freq"`             // 重复频率
	RepeatInterval int             `gorm:"not null;default:1" json:"repeat_interval"`       // 重复间隔（如每2周则设为2）
	NextDate       time.Time       `gorm:"not null" json:"next_date"`                       // 下次执行日期
	EndDate        *time.Time      `json:"end_date"`                                        // 结束日期
	Repetitions    int             `gorm:"default:0" json:"repetitions"`                    // 已执行次数
	MaxRepetitions *int            `json:"max_repetitions"`                                 // 最大执行次数
	IsActive       bool            `gorm:"default:true" json:"is_active"`                   // 是否激活
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`                                  // 软删除时间
}

// TableName 指定表名
func (Recurrence) TableName() string { return "recurrences" }
