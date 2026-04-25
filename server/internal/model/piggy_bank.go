// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PiggyBank 储蓄罐模型，对应 piggy_banks 表
//
// 功能说明：
// - 设定储蓄目标，跟踪储蓄进度
// - 关联到资产账户，从账户中存入/取出资金
// - 支持设定目标日期，追踪储蓄是否按时完成
// - 通过 PiggyEvent 记录每次存取操作
//
// 使用场景：
// - 旅游基金：目标 10000 元，每月存 1000 元
// - 应急储备：目标 30000 元，建立应急资金
// - 大件购买：目标 5000 元，攒钱买新手机
// - 教育基金：目标 50000 元，为子女教育储蓄
//
// 储蓄进度计算：
//   进度 = CurrentAmount / TargetAmount × 100%
//   如：目标 10000 元，已存 6000 元，进度 60%
//
// 与其他模型的关系：
// - Account: 多对一关系，储蓄罐关联到一个资产账户
// - PiggyEvent: 一对多关系，每次存取操作生成一条事件记录
// - Transaction: 间接关联，存取操作可关联到交易
//
// 示例：
//   储蓄罐：旅游基金
//   目标金额：10000 元
//   当前金额：6000 元
//   关联账户：招商银行卡
//   目标日期：2026-12-31
//   进度：60%
type PiggyBank struct {
	// ID 储蓄罐唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的储蓄罐集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询储蓄罐
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Name 储蓄罐名称
	// 用户自定义的储蓄目标名称
	// 如："旅游基金"、"应急储备"、"新手机"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// TargetAmount 目标金额
	// 储蓄的最终目标金额
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	TargetAmount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"target_amount"`

	// CurrentAmount 当前已存金额
	// 累计存入金额 - 累计取出金额
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	// 默认值：0
	CurrentAmount decimal.Decimal `gorm:"type:decimal(19,4);not null;default:0" json:"current_amount"`

	// AccountID 关联账户 ID
	// 资金来源/去向的资产账户
	// 存入时从该账户扣款，取出时退回该账户
	AccountID uint64 `gorm:"not null" json:"account_id"`

	// Order 排序序号
	// 数值越小越靠前，用于控制储蓄罐在列表中的显示顺序
	// 默认值：0
	Order int `gorm:"default:0" json:"order"`

	// TargetDate 目标达成日期（可选）
	// 指针类型表示可为 nil（不设定期限）
	// nil: 不设定期限，按自己节奏存
	// 非 nil: 有明确期限，系统可提醒进度
	TargetDate *time.Time `gorm:"index" json:"target_date,omitempty"`

	// Notes 备注信息
	// 记录储蓄目标的附加说明
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// CreatedAt 储蓄罐创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 储蓄罐最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Account 关联的账户信息
	// 通过 AccountID 外键关联到 accounts 表
	// json:"account,omitempty" 当字段为零值时 JSON 序列化时忽略
	Account Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`

	// Events 存取记录列表
	// 通过 PiggyBankID 外键关联到 piggy_events 表
	// 记录每次存入或取出操作
	// json:"events,omitempty" 当字段为零值时 JSON 序列化时忽略
	Events []PiggyEvent `gorm:"foreignKey:PiggyBankID" json:"events,omitempty"`
}

// TableName 指定 PiggyBank 模型对应的数据库表名为 piggy_banks
func (PiggyBank) TableName() string { return "piggy_banks" }

// PiggyEvent 储蓄事件模型，对应 piggy_events 表
//
// 功能说明：
// - 记录储蓄罐的每次存入或取出操作
// - 支持关联交易记录，实现资金流向追踪
// - 金额正数表示存入，负数表示取出
//
// 使用场景：
// - 存入：每月从工资账户转入 1000 元到旅游基金
// - 取出：从应急储备取出 2000 元支付意外支出
// - 关联交易：存取操作与具体交易关联，便于对账
//
// 示例：
//   存入事件：金额 +1000，关联交易 ID=123，备注"4月储蓄"
//   取出事件：金额 -500，关联交易 ID=456，备注"应急支出"
type PiggyEvent struct {
	// ID 储蓄事件唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// PiggyBankID 储蓄罐 ID
	// 关联到 piggy_banks 表
	// gorm:"index" 创建索引，加速按储蓄罐查询事件记录
	PiggyBankID uint64 `gorm:"not null;index" json:"piggy_bank_id"`

	// Amount 金额
	// 正数：存入金额（如 +1000.00）
	// 负数：取出金额（如 -500.00）
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// TransactionID 关联交易 ID（可选）
	// 指针类型表示可为 nil（手动存取不关联交易）
	// 非 nil: 关联到 transactions 表，表示该存取操作来自一笔交易
	TransactionID *uint64 `gorm:"index" json:"transaction_id,omitempty"`

	// Note 备注信息
	// 记录本次操作的说明，如"4月储蓄"、"应急支出"
	// gorm:"type:text" 支持较长文本
	Note string `gorm:"type:text" json:"note"`

	// CreatedAt 事件创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 PiggyEvent 模型对应的数据库表名为 piggy_events
func (PiggyEvent) TableName() string { return "piggy_events" }
