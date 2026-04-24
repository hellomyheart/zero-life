package request

type CreateTransactionLinkReq struct {
	TransactionID    uint64 `json:"transaction_id" binding:"required"`
	LinkType         string `json:"link_type" binding:"required,oneof=rolled_back reconciled linked"`
	LinkedJournalID  uint64 `json:"linked_journal_id" binding:"required"`
}

type TransactionLinkListReq struct {
	TransactionID *uint64 `form:"transaction_id"`
	Page          int     `form:"page,default=1"`
	PageSize      int     `form:"page_size,default=20"`
}