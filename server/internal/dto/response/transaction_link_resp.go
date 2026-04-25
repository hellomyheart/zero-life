// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// TransactionLinkResp 交易关联响应
type TransactionLinkResp struct {
	ID              uint64    `json:"id"`               // 关联ID
	TransactionID   uint64    `json:"transaction_id"`   // 源交易ID
	LinkType        string    `json:"link_type"`        // 关联类型
	LinkedJournalID uint64    `json:"linked_journal_id"` // 目标交易ID
	CreatedAt       time.Time `json:"created_at"`       // 创建时间
}