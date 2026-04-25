// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateTransactionLinkReq 创建交易关联请求
// 用于在两笔交易之间建立关联关系
type CreateTransactionLinkReq struct {
	TransactionID    uint64 `json:"transaction_id" binding:"required"`                              // 源交易ID
	LinkType         string `json:"link_type" binding:"required,oneof=rolled_back reconciled linked"` // 关联类型：rolled_back(冲正)/reconciled(对账)/linked(关联)
	LinkedJournalID  uint64 `json:"linked_journal_id" binding:"required"`                           // 目标交易ID
}

// TransactionLinkListReq 交易关联列表查询请求
type TransactionLinkListReq struct {
	TransactionID *uint64 `form:"transaction_id"`         // 按交易ID过滤
	Page          int     `form:"page,default=1"`         // 页码
	PageSize      int     `form:"page_size,default=20"`   // 每页数量
}