// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RecurrenceType 循环类型枚举
// 定义循环交易的重复周期，与 Recurrence 模型的 RepeatFreq 类似
// 注意：此模型与 Recurrence 模型功能重叠，后续可能合并
type RecurrenceType string

const (
	// RecurrenceTypeDaily 每日循环
	RecurrenceTypeDaily RecurrenceType = "daily"

	// RecurrenceTypeWeekly 每周循环
	RecurrenceTypeWeekly RecurrenceType = "weekly"

	// RecurrenceTypeMonthly 每月循环
	// 最常见的类型：房租、工资、订阅服务
	RecurrenceTypeMonthly RecurrenceType = "monthly"

	// RecurrenceTypeYearly 每年循环
	// 如：保险费、年检费
	RecurrenceTypeYearly RecurrenceType = "yearly"
)

// RecurringTransaction 循环交易模型，对应 recurring_transactions 表
//
// 功能说明：
// - 定义周期性交易模板，到期自动创建交易
// - 支持多种循环类型（每日/每周/每月/每年）
// - 支持自定义循环间隔（如每2周、每3个月）
// - 支持设定开始和结束日期
// - 通过 RecurringTransactionLog 记录每次生成的交易
//
// 与 Recurrence 模型的区别：
// - Recurrence: 支持更多字段（标签、描述、最大次数等）
// - RecurringTransaction: 更简洁，侧重于基本的循环交易
// - 两者功能重叠，后续可能合并为统一模型
//
// 使用场景：
// - 每月工资：自动创建工资收入交易
// - 每月房租：自动创建房租支出交易
// - 每周保洁：自动创建保洁支出交易
//
// 示例：
//   循环交易：每月房租
//   金额：5000 元
//   循环类型：monthly
//   循环间隔：1（每月一次）
//   开始日期：2026-01-01
//   下次执行：2026-05-01
type RecurringTransaction struct {
	// ID 循环交易唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的循环交易集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Description 交易描述
	// 自动创建交易时的描述文本
	// 如："月度房租"、"工资收入"
	// gorm:"size:500" 限制最大长度为 500 个字符
	Description string `gorm:"not null;size:500" json:"description"`

	// Amount 交易金额
	// 每次自动创建交易时的金额
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// SourceID 源账户 ID
	// 支出交易：资金流出的账户
	// 收入交易：收入来源账户
	SourceID uint64 `gorm:"not null" json:"source_id"`

	// DestinationID 目标账户 ID（可选）
	// 支出交易：支出分类账户
	// 收入交易：资金流入的账户
	// 指针类型表示可为 nil
	DestinationID *uint64 `json:"destination_id"`

	// CategoryID 分类 ID（可选）
	// 自动创建交易时使用的分类
	// 指针类型表示可为 nil
	CategoryID *uint64 `json:"category_id"`

	// Notes 备注信息
	// 记录循环交易的附加说明
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// RecurrenceType 循环类型
	// 取值为 RecurrenceType 枚举：daily/weekly/monthly/yearly
	// gorm:"size:20" 限制最大长度为 20 个字符
	RecurrenceType RecurrenceType `gorm:"not null;size:20" json:"recurrence_type"`

	// RepeatEvery 循环间隔
	// 与 RecurrenceType 配合，实现"每 N 个周期"的效果
	// 如：RecurrenceType=weekly, RepeatEvery=2 表示"每2周"
	// 默认值：1（每个周期）
	RepeatEvery int `gorm:"default:1" json:"repeat_every"`

	// StartDate 开始日期
	// 循环交易的首次执行日期
	StartDate time.Time `gorm:"not null" json:"start_date"`

	// EndDate 结束日期（可选）
	// 指针类型表示可为 nil（无结束日期，永久循环）
	// 非 nil: 到达此日期后不再创建交易
	EndDate *time.Time `json:"end_date,omitempty"`

	// NextOccurrence 下次执行日期
	// 系统在此日期自动创建一笔交易
	// 创建后自动计算下一次执行日期
	NextOccurrence time.Time `gorm:"not null" json:"next_occurrence"`

	// IsActive 是否激活
	// true: 到期自动创建交易
	// false: 暂停自动创建（保留配置）
	// 默认值：true
	IsActive bool `gorm:"default:true" json:"is_active"`

	// CreatedAt 循环交易创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 循环交易最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 RecurringTransaction 模型对应的数据库表名为 recurring_transactions
func (RecurringTransaction) TableName() string { return "recurring_transactions" }

// RecurringTransactionLog 循环交易日志模型，对应 recurring_transaction_logs 表
//
// 功能说明：
// - 记录每次循环交易自动生成的交易
// - 用于追踪循环交易的执行历史
// - 防止重复创建交易（通过唯一性校验）
//
// 使用场景：
// - 查看循环交易的执行记录
// - 排查是否遗漏或重复创建交易
// - 统计循环交易的执行次数
//
// 示例：
//   循环交易 ID=1（每月房租）
//   日志记录：
//     - 2026-01-01 创建交易 ID=101
//     - 2026-02-01 创建交易 ID=156
//     - 2026-03-01 创建交易 ID=210
type RecurringTransactionLog struct {
	// ID 日志记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// RecurringTransactionID 循环交易 ID
	// 关联到 recurring_transactions 表
	// gorm:"index" 创建索引，加速按循环交易查询日志
	RecurringTransactionID uint64 `gorm:"not null;index" json:"recurring_transaction_id"`

	// TransactionID 生成的交易 ID
	// 关联到 transactions 表
	// 记录本次自动创建的交易
	TransactionID uint64 `gorm:"not null" json:"transaction_id"`

	// OccurrenceDate 执行日期
	// 实际创建交易的日期
	OccurrenceDate time.Time `gorm:"not null" json:"occurrence_date"`

	// CreatedAt 日志创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 RecurringTransactionLog 模型对应的数据库表名为 recurring_transaction_logs
func (RecurringTransactionLog) TableName() string { return "recurring_transaction_logs" }
