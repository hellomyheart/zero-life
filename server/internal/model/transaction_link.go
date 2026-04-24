// Package model 数据模型层，对应数据库结构
// TransactionLink 交易关联模型，记录交易之间的关联
package model

import (
	"time"
)

type TransactionLinkType string

const (
	TransactionLinkTypeRelated    TransactionLinkType = "related"
	TransactionLinkTypeReconciled TransactionLinkType = "reconciled"
	TransactionLinkTypeRolledBack TransactionLinkType = "rolled_back"
)

type TransactionJournalLink struct {
	ID              uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID   uint64              `gorm:"not null;index:idx_journal_link" json:"transaction_id"`
	LinkType        TransactionLinkType `gorm:"not null;size:50" json:"link_type"`
	LinkedJournalID uint64              `gorm:"not null;index:idx_journal_link" json:"linked_journal_id"`
	CreatedAt       time.Time           `gorm:"not null" json:"created_at"`
}

func (TransactionJournalLink) TableName() string { return "transaction_journal_links" }