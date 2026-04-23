package request

type CreateTransactionLinkReq struct {
	LinkTypeID    uint64 `json:"link_type_id" binding:"required"`
	SourceID      uint64 `json:"source_id" binding:"required"`
	DestinationID uint64 `json:"destination_id" binding:"required"`
	Comment       string `json:"comment"`
}

type TransactionLinkListReq struct {
	TransactionID uint64 `form:"transaction_id" binding:"required"`
}
