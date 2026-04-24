package request

type CreateTransactionReq struct {
	Type          string            `json:"type" binding:"required,oneof=deposit withdrawal transfer"`
	Date          string            `json:"date" binding:"required"`
	Description   string            `json:"description" binding:"required"`
	Amount        string            `json:"amount" binding:"required"`
	SourceID      uint64            `json:"source_id" binding:"required"`
	DestinationID *uint64           `json:"destination_id"`
	CategoryID    *uint64           `json:"category_id"`
	Notes         string            `json:"notes"`
	Tags          []uint64          `json:"tags"`
	Splits        []CreateSplitReq  `json:"splits"`
}

// CreateSplitReq 创建拆分项请求
type CreateSplitReq struct {
	Amount      string   `json:"amount" binding:"required"`      // 拆分金额
	Description string   `json:"description"`                   // 拆分描述（可选，默认使用父交易描述）
	CategoryID  *uint64  `json:"category_id"`                   // 分类ID（可选）
	Tags        []uint64 `json:"tags"`                          // 标签ID列表（可选）
	Notes       string   `json:"notes"`                         // 备注（可选）
}

type UpdateTransactionReq struct {
	Type          string            `json:"type" binding:"required,oneof=deposit withdrawal transfer"`
	Date          string            `json:"date" binding:"required"`
	Description   string            `json:"description" binding:"required"`
	Amount        string            `json:"amount" binding:"required"`
	SourceID      uint64            `json:"source_id" binding:"required"`
	DestinationID *uint64           `json:"destination_id"`
	CategoryID    *uint64           `json:"category_id"`
	Notes         string            `json:"notes"`
	Tags          []uint64          `json:"tags"`
	Splits        []CreateSplitReq  `json:"splits"`
}

type TransactionListReq struct {
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"page_size,default=20"`
	Type       string  `form:"type"`
	StartDate  string  `form:"start_date"`
	EndDate    string  `form:"end_date"`
	AccountID  *uint64 `form:"account_id"`
	CategoryID *uint64 `form:"category_id"`
	TagID      *uint64 `form:"tag_id"`
	Sort       string  `form:"sort,default=-date"`
}

type TransactionSearchReq struct {
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"page_size,default=20"`
	Keyword    string  `form:"keyword"`                              // 关键词搜索（描述、备注）
	Type       string  `form:"type"`                                 // 交易类型过滤
	StartDate  string  `form:"start_date"`                           // 开始日期
	EndDate    string  `form:"end_date"`                             // 结束日期
	MinAmount  string  `form:"min_amount"`                           // 最小金额
	MaxAmount  string  `form:"max_amount"`                           // 最大金额
	AccountID  *uint64 `form:"account_id"`                           // 账户ID
	CategoryID *uint64 `form:"category_id"`                          // 分类ID
	TagID      *uint64 `form:"tag_id"`                               // 标签ID
	Sort       string  `form:"sort,default=-date"`                   // 排序字段
}

// SplitTransactionReq 拆分交易请求
type SplitTransactionReq struct {
	Splits []CreateSplitReq `json:"splits" binding:"required,min=2"` // 至少拆分为2笔
}
