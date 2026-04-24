// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Reconciliation 对账模型，对应reconciliations表
// 记录账户对账信息，核对账面余额与实际余额
type Reconciliation struct {
	ID               uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64          `gorm:"not null;index" json:"user_id"`                       // 所属用户ID
	AccountID        uint64          `gorm:"not null;index" json:"account_id"`                    // 对账账户ID
	StartDate        time.Time       `gorm:"not null" json:"start_date"`                          // 对账开始日期
	EndDate          time.Time       `gorm:"not null" json:"end_date"`                            // 对账结束日期
	StartBalance     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"start_balance"`    // 期初余额
	EndBalance       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"end_balance"`      // 期末余额
	SubmittedBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"submitted_balance"` // 用户提交的余额
	Difference       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"difference"`       // 差额（提交余额-账面余额）
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (Reconciliation) TableName() string { return "reconciliations" }

// TransactionReconciliation 交易对账模型，对应transaction_reconciliations表
type TransactionReconciliation struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID     uint64    `gorm:"not null;index" json:"account_id"`              // 对账账户ID
	StartDate     time.Time `gorm:"not null" json:"start_date"`                    // 对账开始日期
	EndDate       time.Time `gorm:"not null" json:"end_date"`                      // 对账结束日期
	StartingBalance string  `gorm:"type:decimal(19,4)" json:"starting_balance"`    // 期初余额
	EndingBalance   string  `gorm:"type:decimal(19,4)" json:"ending_balance"`      // 期末余额
	BookBalance   string    `gorm:"type:decimal(19,4)" json:"book_balance"`        // 账面余额
	Difference    string    `gorm:"type:decimal(19,4)" json:"difference"`          // 差额
	Status        string    `gorm:"default:open" json:"status"`                    // 对账状态（open/closed）
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (TransactionReconciliation) TableName() string { return "transaction_reconciliations" }

// ReconciliationEntry 对账条目模型，对应reconciliation_entries表
// 记录对账中每笔交易的匹配情况
type ReconciliationEntry struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReconciliationID     uint64    `gorm:"not null;index" json:"reconciliation_id"`       // 对账记录ID
	TransactionID        uint64    `gorm:"not null;index" json:"transaction_id"`          // 交易ID
	Matched              bool      `gorm:"default:false" json:"matched"`                  // 是否匹配
	AmountDifference     string    `gorm:"type:decimal(19,4)" json:"amount_difference"`   // 金额差异
	CreatedAt            time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (ReconciliationEntry) TableName() string { return "reconciliation_entries" }
