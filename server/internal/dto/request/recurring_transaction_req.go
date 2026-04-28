// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateRecurringTransactionReq 创建循环交易请求
// 用于设置定期重复的交易（如月薪、房租等）
type CreateRecurringTransactionReq struct {
	Description    string  `json:"description" binding:"required"`
	Amount         string  `json:"amount" binding:"required"`
	SourceID       uint64  `json:"source_id" binding:"required"`
	DestinationID  *uint64 `json:"destination_id"`
	CategoryID     *uint64 `json:"category_id"`
	Notes          string  `json:"notes"`
	RecurrenceType string  `json:"recurrence_type" binding:"required,oneof=daily weekly monthly yearly"`
	RepeatEvery    int     `json:"repeat_every"`
	StartDate      string  `json:"start_date" binding:"required"`
	EndDate        *string `json:"end_date"`
	ReminderDays   int     `json:"reminder_days"`
}

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
	IsActive       *bool   `json:"is_active"`
	ReminderDays   *int    `json:"reminder_days"`
}

// RecurringTransactionListReq 循环交易列表查询请求
type RecurringTransactionListReq struct {
	Page     int    `form:"page,default=1"`        // 页码
	PageSize int    `form:"page_size,default=20"`  // 每页数量
	Active   *bool  `form:"active"`                // 按启用状态过滤
}