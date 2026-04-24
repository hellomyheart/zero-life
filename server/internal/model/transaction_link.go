// Package model 数据模型层，对应数据库结构
package model

import (
	"time"
)

// TransactionLinkType 交易关联类型枚举
// 定义交易之间关联关系的类型
type TransactionLinkType string

const (
	// TransactionLinkTypeRelated 关联：两笔交易有逻辑关联（如分期付款的各期）
	TransactionLinkTypeRelated TransactionLinkType = "related"

	// TransactionLinkTypeReconciled 已对账：两笔交易通过对账确认匹配
	TransactionLinkTypeReconciled TransactionLinkType = "reconciled"

	// TransactionLinkTypeRolledBack 回滚：一笔交易是另一笔的回滚/撤销
	TransactionLinkTypeRolledBack TransactionLinkType = "rolled_back"
)

// TransactionJournalLink 交易日志关联模型，对应transaction_journal_links表
// 记录两笔交易之间的关联关系，支持多种关联类型
// 例如：一笔退款交易关联到原始消费交易
type TransactionJournalLink struct {
	// ID 关联记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// TransactionID 源交易ID，关联关系的发起方
	// gorm:"index:idx_journal_link" 创建复合索引的一部分，加速关联查询
	TransactionID uint64 `gorm:"not null;index:idx_journal_link" json:"transaction_id"`

	// LinkType 关联类型，取值为TransactionLinkType枚举
	// gorm:"size:50" 限制类型名称最大长度为50个字符
	LinkType TransactionLinkType `gorm:"not null;size:50" json:"link_type"`

	// LinkedJournalID 目标交易日志ID，关联关系的接收方
	// gorm:"index:idx_journal_link" 与TransactionID组成复合索引
	// 复合索引可高效查询"某笔交易的所有关联"以及"两笔交易是否已关联"
	LinkedJournalID uint64 `gorm:"not null;index:idx_journal_link" json:"linked_journal_id"`

	// CreatedAt 关联创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定数据库表名为transaction_journal_links
func (TransactionJournalLink) TableName() string { return "transaction_journal_links" }
