// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionType 交易类型枚举
// 定义系统中支持的三种交易类型，对应复式记账法的基本操作
type TransactionType string

const (
	// TransactionTypeDeposit 存款/收入：资金流入
	// 源账户为收入账户，目标账户为资产账户
	// 如：工资到账、红包收入
	TransactionTypeDeposit TransactionType = "deposit"

	// TransactionTypeWithdrawal 取款/支出：资金流出
	// 源账户为资产账户，目标账户为支出账户
	// 如：餐饮消费、购物支出
	TransactionTypeWithdrawal TransactionType = "withdrawal"

	// TransactionTypeTransfer 转账：账户间转移
	// 源账户和目标账户都为资产账户
	// 如：银行卡转支付宝、信用卡还款
	TransactionTypeTransfer TransactionType = "transfer"
)

// Transaction 交易模型，对应 transactions 表
//
// 功能说明：
// - 记录所有的财务交易，是系统的核心模型
// - 支持三种交易类型：收入（deposit）、支出（withdrawal）、转账（transfer）
// - 支持交易拆分：一笔交易可拆分为多笔子交易
// - 支持标签、分类、账单关联等高级功能
// - 支持对账标记
//
// 复式记账原理：
// - 每笔交易涉及两个账户（源账户和目标账户）
// - 支出：资产账户（源）-> 支出账户（目标）
// - 收入：收入账户（源）-> 资产账户（目标）
// - 转账：资产账户A（源）-> 资产账户B（目标）
//
// 交易拆分：
// - 一笔交易可以拆分为多笔子交易（通过 ParentID 关联）
// - 子交易共享父交易的日期和描述
// - 子交易可以有不同的分类、金额和标签
// - 如：超市购物 200 元 = 食品 120 元 + 日用品 80 元
//
// 与其他模型的关系：
// - Account: 通过 SourceID/DestinationID 关联源/目标账户
// - Category: 通过 CategoryID 关联分类
// - Tag: 多对多关系，通过 transaction_tags 中间表
// - Bill: 通过 BillID 关联账单
// - Transaction: 自关联，通过 ParentID 实现交易拆分
//
// 示例：
//   支出交易：星巴克咖啡
//   类型：withdrawal
//   金额：35 元
//   源账户：招商银行卡（资产账户）
//   目标账户：餐饮（支出账户）
//   分类：咖啡
//   标签：["网购", "待报销"]
type Transaction struct {
	// ID 交易唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的交易集合，用户间数据隔离
	// gorm:"index:idx_user_date" 与 Date 组成复合索引，加速按用户+日期查询
	UserID uint64 `gorm:"not null;index:idx_user_date" json:"user_id"`

	// Type 交易类型
	// 取值为 TransactionType 枚举：deposit/withdrawal/transfer
	// 决定交易的账户流向
	// gorm:"size:20" 限制最大长度为 20 个字符
	Type TransactionType `gorm:"not null;size:20" json:"type"`

	// Date 交易日期
	// 记录交易发生的日期，用于报表统计
	// gorm:"index:idx_user_date" 与 UserID 组成复合索引
	Date time.Time `gorm:"not null;index:idx_user_date" json:"date"`

	// Description 交易描述
	// 交易的简要说明，如"星巴克咖啡"、"工资"
	// gorm:"size:500" 限制最大长度为 500 个字符
	Description string `gorm:"not null;size:500" json:"description"`

	// Amount 交易金额
	// 使用 decimal 类型确保金额精度，避免浮点数误差
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// SourceID 源账户 ID
	// 支出交易：资产账户（资金流出）
	// 收入交易：收入账户（收入来源）
	// 转账交易：转出账户
	// gorm:"index:idx_user_source" 与 UserID 组成复合索引，加速按用户+源账户查询
	SourceID uint64 `gorm:"not null;index:idx_user_source" json:"source_id"`

	// DestinationID 目标账户 ID（可选）
	// 支出交易：支出账户（资金去向）
	// 收入交易：资产账户（资金流入）
	// 转账交易：转入账户
	// 指针类型表示可为 nil（某些简单交易可能没有目标账户）
	DestinationID *uint64 `json:"destination_id"`

	// CategoryID 分类 ID（可选）
	// 关联到 categories 表，用于分类统计和预算控制
	// 指针类型表示可为 nil（交易可以暂不分类）
	// gorm:"index:idx_user_category" 与 UserID 组成复合索引
	CategoryID *uint64 `gorm:"index:idx_user_category" json:"category_id"`

	// Notes 备注信息
	// 记录交易的附加说明
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// BillID 关联账单 ID（可选）
	// 关联到 bills 表，表示该交易对应哪个账单
	// 指针类型表示可为 nil（非账单关联的交易）
	BillID *uint64 `json:"bill_id"`

	// ParentID 父交易 ID（用于交易拆分）
	// nil: 这是一笔独立交易（或父交易）
	// 非 nil: 这是父交易的子交易（拆分项）
	// gorm:"index" 创建索引，加速查询子交易
	ParentID *uint64 `gorm:"index" json:"parent_id"`

	// IsReconciled 是否已对账
	// true: 该交易已与银行对账单核对
	// false: 尚未核对
	// 默认值：false
	IsReconciled bool `gorm:"default:false" json:"is_reconciled"`

	// CreatedAt 交易创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 交易最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Source 源账户信息
	// 通过 SourceID 外键关联到 accounts 表
	// json:"source,omitempty" 当字段为零值时 JSON 序列化时忽略
	Source Account `gorm:"foreignKey:SourceID" json:"source,omitempty"`

	// Destination 目标账户信息（可选）
	// 通过 DestinationID 外键关联到 accounts 表
	// 指针类型，因为目标账户可能为空
	// json:"destination,omitempty" 当字段为零值时 JSON 序列化时忽略
	Destination *Account `gorm:"foreignKey:DestinationID" json:"destination,omitempty"`

	// Category 分类信息（可选）
	// 通过 CategoryID 外键关联到 categories 表
	// 指针类型，因为分类可能为空
	// json:"category,omitempty" 当字段为零值时 JSON 序列化时忽略
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// Tags 标签列表（多对多关系）
	// 通过 transaction_tags 中间表实现
	// 一个交易可以有多个标签
	// json:"tags,omitempty" 当字段为零值时 JSON 序列化时忽略
	Tags []Tag `gorm:"many2many:transaction_tags" json:"tags,omitempty"`

	// Splits 拆分交易列表
	// 通过 ParentID 外键关联，获取当前交易的所有子交易
	// 仅父交易有子交易，子交易的 Splits 为空
	// json:"splits,omitempty" 当字段为零值时 JSON 序列化时忽略
	Splits []Transaction `gorm:"foreignKey:ParentID" json:"splits,omitempty"`
}

// TableName 指定 Transaction 模型对应的数据库表名为 transactions
func (Transaction) TableName() string { return "transactions" }
