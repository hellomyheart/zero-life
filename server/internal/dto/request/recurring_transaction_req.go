// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateRecurringTransactionReq 创建循环交易请求
// 用于设置定期重复的交易（如月薪、房租等）
type CreateRecurringTransactionReq struct {
	Description    string  `json:"description" binding:"required"`                              // 交易描述
	Amount         string  `json:"amount" binding:"required"`                                   // 交易金额
	SourceID       uint64  `json:"source_id" binding:"required"`                                // 源账户ID
	DestinationID  *uint64 `json:"destination_id"`                                              // 目标账户ID
	CategoryID     *uint64 `json:"category_id"`                                                 // 分类ID
	Notes          string  `json:"notes"`                                                       // 备注
	RecurrenceType string  `json:"recurrence_type" binding:"required,oneof=daily weekly monthly yearly"` // 重复类型
	RepeatEvery    int     `json:"repeat_every"`                                                // 重复间隔
	StartDate      string  `json:"start_date" binding:"required"`                               // 开始日期
	EndDate        *string `json:"end_date"`                                                    // 结束日期（可选）
}

// UpdateRecurringTransactionReq 更新循环交易请求
type UpdateRecurringTransactionReq struct {
	Description    string  `json:"description"`
	Amount         string  `json:"amount"`
	SourceID       *uint64 `json:"source_id"`
	DestinationID  *uint64 `json:"destination_id"`
	CategoryID     *uint64 `json:"category_id"`
	Notes          string  `json:"notes"`
	RecurrenceType string  `json:"recurrence_type" binding:"omitempty,oneof=daily weekly monthly yearly"`
	RepeatEvery    *int    `json:"repeat_every"`
	StartDate      string  `json:"start_date"`
	EndDate        *string `json:"end_date"`
	IsActive       *bool   `json:"is_active"` // 是否启用
}

// RecurringTransactionListReq 循环交易列表查询请求
type RecurringTransactionListReq struct {
	Page     int    `form:"page,default=1"`        // 页码
	PageSize int    `form:"page_size,default=20"`  // 每页数量
	Active   *bool  `form:"active"`                // 按启用状态过滤
}