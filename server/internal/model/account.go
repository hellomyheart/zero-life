// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccountType 账户类型枚举
// 定义系统中支持的四种账户类型，采用复式记账法的账户分类
type AccountType string

const (
	// AccountTypeAsset 资产账户：银行存款、现金、投资等
	// 资产账户的余额代表用户拥有的财富
	// 例如：招商银行卡、支付宝余额、微信零钱
	AccountTypeAsset AccountType = "asset"

	// AccountTypeExpense 支出账户：用于记录各类支出
	// 支出账户的余额代表累计支出金额
	// 例如：餐饮、交通、购物、娱乐
	AccountTypeExpense AccountType = "expense"

	// AccountTypeRevenue 收入账户：用于记录各类收入
	// 收入账户的余额代表累计收入金额
	// 例如：工资、奖金、投资收益
	AccountTypeRevenue AccountType = "revenue"

	// AccountTypeLiability 负债账户：信用卡、贷款等
	// 负债账户的余额代表用户欠款的金额
	// 例如：招商信用卡、房贷、车贷
	AccountTypeLiability AccountType = "liability"
)

// Account 账户模型，对应 accounts 表
//
// 功能说明：
// - 管理用户的各类财务账户，采用复式记账法分类
// - 支持四种账户类型：资产、支出、收入、负债
// - 支持虚拟账户（如预算跟踪账户）
// - 自动计算当前余额
//
// 复式记账原理：
// - 每笔交易涉及至少两个账户（一借一贷）
// - 支出交易：资产账户（减少）-> 支出账户（增加）
// - 收入交易：收入账户（增加）-> 资产账户（增加）
// - 转账交易：资产账户A（减少）-> 资产账户B（增加）
//
// 与其他模型的关系：
// - Currency: 每个账户关联一种货币
// - Transaction: 交易通过 SourceID/DestinationID 关联账户
// - PiggyBank: 储蓄罐关联到资产账户
// - Reconciliation: 对账记录关联到资产账户
//
// 示例：
//   资产账户：招商银行卡 - CNY - 初始余额 10000 - 当前余额 8500
//   支出账户：餐饮 - CNY - 初始余额 0 - 当前余额 3200（累计支出）
//   收入账户：工资 - CNY - 初始余额 0 - 当前余额 30000（累计收入）
//   负债账户：招商信用卡 - CNY - 初始余额 0 - 当前余额 5000（欠款）
type Account struct {
	// ID 账户唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的账户集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询账户
	UserID uint64 `gorm:"not null;index;uniqueIndex:idx_user_account_num_del" json:"user_id"`

	// Name 账户名称
	// 用户自定义的账户名称，便于识别
	// 如："招商银行卡"、"支付宝"、"餐饮支出"
	// gorm:"size:255" 限制最大长度为 255 个字符
	Name string `gorm:"not null;size:255" json:"name"`

	// AccountNumber 账户号
	// 用户自定义的账户编号，如银行卡号、信用卡号等
	// 同一用户内唯一（含软删除），创建后不可修改
	// gorm:"not null;size:100" 限制最大长度为 100 个字符
	// 复合唯一索引 idx_user_account_num_del: (user_id, account_number, deleted_at)
	AccountNumber string `gorm:"not null;size:100;uniqueIndex:idx_user_account_num_del" json:"account_number"`

	// Type 账户类型
	// 取值为 AccountType 枚举：asset/expense/revenue/liability
	// 决定账户在复式记账中的角色
	// gorm:"size:20" 限制最大长度为 20 个字符
	Type AccountType `gorm:"not null;size:20" json:"type"`

	// CurrencyID 货币 ID
	// 关联到 currencies 表，指定账户使用的货币
	// 账户余额使用该货币的单位
	// 不同货币的账户余额不能直接相加，需要汇率折算
	CurrencyID uint64 `gorm:"not null" json:"currency_id"`

	// InitialBalance 初始余额
	// 账户创建时的起始余额
	// 使用 decimal 类型确保金额精度，避免浮点数误差
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	// 默认值：0
	InitialBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"initial_balance"`

	// CurrentBalance 当前余额（系统自动计算）
	// = 初始余额 + 所有收入 - 所有支出
	// 资产账户：正值表示余额，负值表示透支
	// 负债账户：正值表示欠款金额
	// 支出账户：正值表示累计支出
	// 收入账户：正值表示累计收入
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	// 默认值：0
	CurrentBalance decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"current_balance"`

	// IsVirtual 是否为虚拟账户
	// true: 虚拟账户，不对应真实资金（如预算跟踪账户）
	// false: 真实账户，对应实际资金（如银行账户）
	// 默认值：false
	IsVirtual bool `gorm:"default:false" json:"is_virtual"`

	// Notes 备注信息
	// 记录账户的附加说明，如账户用途、开户行信息等
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// CreatedAt 账户创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 账户最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_user_account_num_del" json:"-"`

	// Currency 关联的货币信息
	// 通过 CurrencyID 外键关联到 currencies 表
	// json:"currency,omitempty" 当字段为零值时 JSON 序列化时忽略
	Currency Currency `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"`
}

// TableName 指定 Account 模型对应的数据库表名为 accounts
func (Account) TableName() string { return "accounts" }
