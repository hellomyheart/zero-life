// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RepeatFreq 重复频率枚举
// 定义定期交易的重复周期
type RepeatFreq string

const (
	// RepeatFreqDaily 每日重复
	// 如：每日午餐支出
	RepeatFreqDaily RepeatFreq = "daily"

	// RepeatFreqWeekly 每周重复
	// 如：每周保洁费用
	RepeatFreqWeekly RepeatFreq = "weekly"

	// RepeatFreqMonthly 每月重复
	// 最常见的频率：房租、工资、订阅服务
	RepeatFreqMonthly RepeatFreq = "monthly"

	// RepeatFreqYearly 每年重复
	// 如：保险费、年检费、会员费
	RepeatFreqYearly RepeatFreq = "yearly"
)

// Recurrence 定期交易模型，对应 recurrences 表
//
// 功能说明：
// - 定义周期性交易模板，到期自动创建交易
// - 支持多种重复频率（每日/每周/每月/每年）
// - 支持自定义重复间隔（如每2周、每3个月）
// - 支持设定结束条件（按日期或按次数）
// - 支持关联分类、标签等交易属性
//
// 使用场景：
// - 工资收入：每月 15 日自动创建工资收入交易
// - 房租支出：每月 1 日自动创建房租支出交易
// - 订阅服务：每月自动创建 Netflix/Spotify 支出
// - 保险费用：每年自动创建保险支出交易
//
// 与 Bill 的区别：
// - Bill: 侧重于账单到期提醒，是"被动"的
// - Recurrence: 侧重于自动创建交易，是"主动"的
//
// 结束条件：
// - EndDate: 到达指定日期后停止
// - MaxRepetitions: 达到指定次数后停止
// - 两者都不设则永久重复
//
// 示例：
//   定期交易：每月工资
//   金额：15000 元
//   频率：每月（RepeatFreq=monthly, RepeatInterval=1）
//   源账户：收入账户
//   目标账户：招商银行卡
//   下次执行：2026-05-15
type Recurrence struct {
	// ID 定期交易唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的定期交易集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Title 定期交易名称
	// 如："每月工资"、"房租"、"Netflix 订阅"
	// gorm:"size:255" 限制最大长度为 255 个字符
	Title string `gorm:"not null;size:255" json:"title"`

	// Type 交易类型
	// 取值为 TransactionType 枚举：deposit/withdrawal/transfer
	// deposit: 收入类定期交易（如工资）
	// withdrawal: 支出类定期交易（如房租）
	// transfer: 转账类定期交易（如定期转账到储蓄账户）
	// gorm:"size:20" 限制最大长度为 20 个字符
	Type TransactionType `gorm:"not null;size:20" json:"type"`

	// Amount 交易金额
	// 每次自动创建交易时的金额
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// SourceID 源账户 ID
	// 支出交易：资金流出的账户（如银行卡）
	// 收入交易：收入来源账户（如工资收入账户）
	// 转账交易：转出账户
	SourceID uint64 `gorm:"not null" json:"source_id"`

	// DestinationID 目标账户 ID（可选）
	// 支出交易：支出分类账户
	// 收入交易：资金流入的账户（如银行卡）
	// 转账交易：转入账户
	// 指针类型表示可为 nil
	DestinationID *uint64 `gorm:"index" json:"destination_id"`

	// CategoryID 分类 ID（可选）
	// 自动创建交易时使用的分类
	// 如：房租定期交易关联"居住"分类
	// 指针类型表示可为 nil
	CategoryID *uint64 `gorm:"index" json:"category_id"`

	// Description 交易描述
	// 自动创建交易时的描述文本
	// gorm:"size:500" 限制最大长度为 500 个字符
	Description string `gorm:"size:500" json:"description"`

	// Notes 备注信息
	// 记录定期交易的附加说明
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// TagNames 标签名称（逗号分隔）
	// 自动创建交易时添加的标签
	// 如："固定支出,房租"
	// gorm:"size:500" 限制最大长度为 500 个字符
	TagNames string `gorm:"size:500" json:"tag_names"`

	// RepeatFreq 重复频率
	// 取值为 RepeatFreq 枚举：daily/weekly/monthly/yearly
	// 与 RepeatInterval 配合使用
	// gorm:"size:20" 限制最大长度为 20 个字符
	RepeatFreq RepeatFreq `gorm:"not null;size:20" json:"repeat_freq"`

	// RepeatInterval 重复间隔
	// 与 RepeatFreq 配合，实现"每 N 个周期"的效果
	// 如：RepeatFreq=weekly, RepeatInterval=2 表示"每2周"
	// 如：RepeatFreq=monthly, RepeatInterval=3 表示"每3个月"
	// 默认值：1（每个周期）
	RepeatInterval int `gorm:"not null;default:1" json:"repeat_interval"`

	// NextDate 下次执行日期
	// 系统在此日期自动创建一笔交易
	// 创建后自动计算下一次执行日期
	NextDate time.Time `gorm:"not null" json:"next_date"`

	// EndDate 结束日期（可选）
	// 指针类型表示可为 nil（无结束日期，永久重复）
	// 非 nil: 到达此日期后不再创建交易
	EndDate *time.Time `gorm:"index" json:"end_date"`

	// Repetitions 已执行次数
	// 记录该定期交易已经自动创建了多少笔交易
	// 默认值：0
	Repetitions int `gorm:"default:0" json:"repetitions"`

	// MaxRepetitions 最大执行次数（可选）
	// 指针类型表示可为 nil（无次数限制）
	// 非 nil: 达到指定次数后不再创建交易
	MaxRepetitions *int `gorm:"type:integer" json:"max_repetitions"`

	// IsActive 是否激活
	// true: 到期自动创建交易
	// false: 暂停自动创建（保留配置）
	// 默认值：true
	IsActive bool `gorm:"default:true" json:"is_active"`

	// CreatedAt 定期交易创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 定期交易最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 Recurrence 模型对应的数据库表名为 recurrences
func (Recurrence) TableName() string { return "recurrences" }
