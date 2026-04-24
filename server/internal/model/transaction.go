package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionType 交易类型枚举
type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "deposit"    // 存款/收入：资金流入
	TransactionTypeWithdrawal  TransactionType = "withdrawal" // 取款/支出：资金流出
	TransactionTypeTransfer    TransactionType = "transfer"   // 转账：账户间转移
)

// Transaction 交易模型
// 记录所有的财务交易，包括收入、支出和转账
// 支持交易拆分、标签、分类等高级功能
type Transaction struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`                  // 交易ID，主键自增
	UserID        uint64          `gorm:"not null;index:idx_user_date" json:"user_id"`         // 所属用户ID
	Type          TransactionType `gorm:"not null;size:20" json:"type"`                        // 交易类型（deposit/withdrawal/transfer）
	Date          time.Time       `gorm:"not null;index:idx_user_date" json:"date"`            // 交易日期
	Description   string          `gorm:"not null;size:500" json:"description"`                // 交易描述
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`           // 交易金额
	SourceID      uint64          `gorm:"not null;index:idx_user_source" json:"source_id"`     // 源账户ID（支出账户或转出账户）
	DestinationID *uint64         `json:"destination_id"`                                      // 目标账户ID（收入账户或转入账户）
	CategoryID    *uint64         `gorm:"index:idx_user_category" json:"category_id"`          // 分类ID
	Notes         string          `gorm:"type:text" json:"notes"`                              // 备注信息
	BillID        *uint64         `json:"bill_id"`                                             // 关联账单ID
	ParentID      *uint64         `gorm:"index" json:"parent_id"`                              // 父交易ID（用于交易拆分）
	IsReconciled  bool            `gorm:"default:false" json:"is_reconciled"`                  // 是否已对账
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`                          // 创建时间
	UpdatedAt     time.Time       `gorm:"not null" json:"updated_at"`                          // 更新时间
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`                                      // 软删除时间

	Source      Account       `gorm:"foreignKey:SourceID" json:"source,omitempty"`         // 源账户信息
	Destination *Account      `gorm:"foreignKey:DestinationID" json:"destination,omitempty"` // 目标账户信息
	Category    *Category     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`     // 分类信息
	Tags        []Tag         `gorm:"many2many:transaction_tags" json:"tags,omitempty"`    // 标签列表
	Splits      []Transaction `gorm:"foreignKey:ParentID" json:"splits,omitempty"`         // 拆分交易列表
}

// TableName 指定表名
func (Transaction) TableName() string { return "transactions" }
