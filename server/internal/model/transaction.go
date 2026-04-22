package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "deposit"
	TransactionTypeWithdrawal  TransactionType = "withdrawal"
	TransactionTypeTransfer    TransactionType = "transfer"
)

type Transaction struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64          `gorm:"not null;index:idx_user_date" json:"user_id"`
	Type          TransactionType `gorm:"not null;size:20" json:"type"`
	Date          time.Time       `gorm:"not null;index:idx_user_date" json:"date"`
	Description   string          `gorm:"not null;size:500" json:"description"`
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	SourceID      uint64          `gorm:"not null;index:idx_user_source" json:"source_id"`
	DestinationID *uint64         `json:"destination_id"`
	CategoryID    *uint64         `gorm:"index:idx_user_category" json:"category_id"`
	Notes         string          `gorm:"type:text" json:"notes"`
	BillID        *uint64         `json:"bill_id"`
	ParentID      *uint64         `gorm:"index" json:"parent_id"`
	CreatedAt     time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`

	Source      Account  `gorm:"foreignKey:SourceID" json:"source,omitempty"`
	Destination *Account `gorm:"foreignKey:DestinationID" json:"destination,omitempty"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag    `gorm:"many2many:transaction_tags" json:"tags,omitempty"`
	Splits      []Transaction `gorm:"foreignKey:ParentID" json:"splits,omitempty"`
}

func (Transaction) TableName() string { return "transactions" }
