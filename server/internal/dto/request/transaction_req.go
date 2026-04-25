// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateTransactionReq 创建交易请求
// 用于创建新的交易记录，支持存款、取款和转账三种类型
type CreateTransactionReq struct {
	Type          string            `json:"type" binding:"required,oneof=deposit withdrawal transfer"` // 交易类型：deposit(存款)/withdrawal(取款)/transfer(转账)
	Date          string            `json:"date" binding:"required"`                                   // 交易日期
	Description   string            `json:"description" binding:"required"`                            // 交易描述
	Amount        string            `json:"amount" binding:"required"`                                 // 交易金额，必须大于0
	SourceID      uint64            `json:"source_id" binding:"required"`                              // 源账户ID（支出账户）
	DestinationID *uint64           `json:"destination_id"`                                            // 目标账户ID（转账时必填）
	CategoryID    *uint64           `json:"category_id"`                                               // 分类ID
	Notes         string            `json:"notes"`                                                     // 备注
	Tags          []uint64          `json:"tags"`                                                      // 标签ID列表
	Splits        []CreateSplitReq  `json:"splits"`                                                    // 拆分项列表
}

// CreateSplitReq 创建拆分项请求
type CreateSplitReq struct {
	Amount      string   `json:"amount" binding:"required"`      // 拆分金额
	Description string   `json:"description"`                   // 拆分描述（可选，默认使用父交易描述）
	CategoryID  *uint64  `json:"category_id"`                   // 分类ID（可选）
	Tags        []uint64 `json:"tags"`                          // 标签ID列表（可选）
	Notes       string   `json:"notes"`                         // 备注（可选）
}

// UpdateTransactionReq 更新交易请求
// 用于修改已有交易的信息，字段含义与CreateTransactionReq相同
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

// TransactionListReq 交易列表查询请求
// 用于分页查询和过滤交易列表
type TransactionListReq struct {
	Page       int     `form:"page,default=1"`        // 页码
	PageSize   int     `form:"page_size,default=20"`  // 每页数量
	Type       string  `form:"type"`                  // 按交易类型过滤
	StartDate  string  `form:"start_date"`            // 开始日期过滤
	EndDate    string  `form:"end_date"`              // 结束日期过滤
	AccountID  *uint64 `form:"account_id"`            // 按账户ID过滤
	CategoryID *uint64 `form:"category_id"`           // 按分类ID过滤
	TagID      *uint64 `form:"tag_id"`                // 按标签ID过滤
	Sort       string  `form:"sort,default=-date"`    // 排序字段，默认按日期倒序
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
