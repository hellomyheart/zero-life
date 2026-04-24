package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Reconciliation struct {
	ID               uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64          `gorm:"not null;index" json:"user_id"`
	AccountID        uint64          `gorm:"not null;index" json:"account_id"`
	StartDate        time.Time       `gorm:"not null" json:"start_date"`
	EndDate          time.Time       `gorm:"not null" json:"end_date"`
	StartBalance     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"start_balance"`
	EndBalance       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"end_balance"`
	SubmittedBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"submitted_balance"`
	Difference       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"difference"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
}

func (Reconciliation) TableName() string { return "reconciliations" }

type TransactionReconciliation struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID     uint64    `gorm:"not null;index" json:"account_id"`
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	StartingBalance string  `gorm:"type:decimal(19,4)" json:"starting_balance"`
	EndingBalance   string  `gorm:"type:decimal(19,4)" json:"ending_balance"`
	BookBalance   string    `gorm:"type:decimal(19,4)" json:"book_balance"`
	Difference    string    `gorm:"type:decimal(19,4)" json:"difference"`
	Status        string    `gorm:"default:open" json:"status"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
}

func (TransactionReconciliation) TableName() string { return "transaction_reconciliations" }

type ReconciliationEntry struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReconciliationID     uint64    `gorm:"not null;index" json:"reconciliation_id"`
	TransactionID        uint64    `gorm:"not null;index" json:"transaction_id"`
	Matched              bool      `gorm:"default:false" json:"matched"`
	AmountDifference     string    `gorm:"type:decimal(19,4)" json:"amount_difference"`
	CreatedAt            time.Time `gorm:"not null" json:"created_at"`
}

func (ReconciliationEntry) TableName() string { return "reconciliation_entries" }