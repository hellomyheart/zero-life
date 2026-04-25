// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// BulkEditReq 批量编辑交易请求
// 批量修改交易的分类、备注和标签
type BulkEditReq struct {
	IDs        []uint64 `json:"ids" binding:"required,min=1"` // 交易ID列表，至少1个
	CategoryID *uint64  `json:"category_id"`                  // 统一设置分类ID
	Notes      string   `json:"notes"`                        // 统一设置备注
	TagIDs     []uint64 `json:"tag_ids"`                      // 统一添加标签ID列表
}

// BulkDeleteReq 批量删除交易请求
type BulkDeleteReq struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"` // 交易ID列表，至少1个
}

// ConvertReq 交易类型转换请求
// 将交易从一种类型转换为另一种类型
type ConvertReq struct {
	Type          string  `json:"type" binding:"required,oneof=deposit withdrawal transfer"` // 目标交易类型
	SourceID      *uint64 `json:"source_id"`                                                 // 源账户ID
	DestinationID *uint64 `json:"destination_id"`                                            // 目标账户ID
}
