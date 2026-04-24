// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Currency 货币模型，对应 currencies 表
// 
// 功能说明：
// - 管理系统支持的货币类型
// - 定义货币的基本属性（代码、名称、符号、小数位）
// - 支持启用/禁用货币
// - 支持设置默认货币
//
// 使用场景：
// - 多货币账户：如美元账户、欧元账户
// - 多货币交易：如海外消费、跨境转账
// - 报表折算：将所有货币折算为本位币
//
// 常见货币示例：
//   CNY - 人民币 - ¥ - 2 位小数
//   USD - 美元 - $ - 2 位小数
//   EUR - 欧元 - € - 2 位小数
//   JPY - 日元 - ¥ - 0 位小数
//   GBP - 英镑 - £ - 2 位小数
//
// 与账户的关系：
// - 每个账户必须关联一种货币
// - 账户的余额使用该货币的单位
// - 不同货币的账户余额不能直接相加，需要汇率折算
type Currency struct {
	// ID 货币唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// Code 货币代码（ISO 4217 标准）
	// 3 个大写字母，如：CNY, USD, EUR
	// gorm:"uniqueIndex" 创建唯一索引，确保代码不重复
	// gorm:"size:3" 固定 3 个字符
	Code string `gorm:"uniqueIndex;not null;size:3" json:"code"`
	
	// Name 货币名称
	// 如："人民币"、"美元"、"欧元"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`
	
	// Symbol 货币符号
	// 用于前端显示，如：¥, $, €, £
	// gorm:"size:10" 限制最大长度为 10 个字符
	Symbol string `gorm:"not null;size:10" json:"symbol"`
	
	// DecimalPlaces 小数位数
	// 大多数货币为 2 位（如 CNY, USD, EUR）
	// 日元（JPY）为 0 位
	// 某些货币可能为 3 位（如 BHD）
	// 默认值：2
	DecimalPlaces int `gorm:"default:2" json:"decimal_places"`
	
	// IsEnabled 是否启用
	// true: 可以在创建账户/交易时选择
	// false: 隐藏，不可选择（但历史数据保留）
	// 默认值：true
	IsEnabled bool `gorm:"default:true" json:"is_enabled"`
	
	// IsDefault 是否为默认货币
	// true: 系统默认使用此货币
	// false: 非默认货币
	// 注意：系统中只能有一个默认货币
	// 默认值：false
	IsDefault bool `gorm:"default:false" json:"is_default"`
	
	// CreatedAt 货币创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 货币最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定 Currency 模型对应的数据库表名为 currencies
func (Currency) TableName() string { return "currencies" }

// ExchangeRate 汇率模型，对应 exchange_rates 表
// 
// 功能说明：
// - 记录不同货币之间的兑换比率
// - 支持汇率历史查询（通过 UpdatedAt）
// - 用于多货币报表折算
//
// 使用场景：
// - 账户余额折算：将所有货币折算为本位币
// - 交易金额折算：计算交易的本位币金额
// - 净值计算：统一货币单位后计算总资产
//
// 汇率表示：
//   FromCurrencyID = 1 (CNY)
//   ToCurrencyID = 2 (USD)
//   Rate = 0.14
//   含义：1 CNY = 0.14 USD
//
// 折算公式：
//   目标金额 = 源金额 × Rate
//   100 CNY × 0.14 = 14 USD
//
// 汇率更新策略：
// - 可以手动更新汇率
// - 可以对接外部 API 自动更新（如央行汇率）
// - 建议每天更新一次
// - 历史汇率用于历史报表的准确折算
type ExchangeRate struct {
	// ID 汇率记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// FromCurrencyID 源货币 ID
	// 关联到 currencies 表的 ID
	// gorm:"index:idx_currency_pair,unique" 与 ToCurrencyID 组成复合唯一索引
	// 确保同一货币对只有一条记录
	FromCurrencyID uint64 `gorm:"not null;index:idx_currency_pair,unique" json:"from_currency_id"`
	
	// ToCurrencyID 目标货币 ID
	// 关联到 currencies 表的 ID
	// gorm:"index:idx_currency_pair,unique" 与 FromCurrencyID 组成复合唯一索引
	ToCurrencyID uint64 `gorm:"not null;index:idx_currency_pair,unique" json:"to_currency_id"`
	
	// Rate 兑换比率
	// 使用 decimal 类型确保精度
	// gorm:"type:decimal(19,8)" 支持最高 8 位小数，确保汇率精度
	Rate decimal.Decimal `gorm:"type:decimal(19,8);not null" json:"rate"`
	
	// UpdatedAt 汇率最后更新时间
	// 用于判断汇率是否过期
	// 建议超过 7 天的汇率需要更新
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定 ExchangeRate 模型对应的数据库表名为 exchange_rates
func (ExchangeRate) TableName() string { return "exchange_rates" }
