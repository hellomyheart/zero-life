package request

type BulkEditReq struct {
	IDs        []uint64 `json:"ids" binding:"required,min=1"`
	CategoryID *uint64  `json:"category_id"`
	Notes      string   `json:"notes"`
	TagIDs     []uint64 `json:"tag_ids"`
}

type BulkDeleteReq struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

type ConvertReq struct {
	Type          string  `json:"type" binding:"required,oneof=deposit withdrawal transfer"`
	SourceID      *uint64 `json:"source_id"`
	DestinationID *uint64 `json:"destination_id"`
}
