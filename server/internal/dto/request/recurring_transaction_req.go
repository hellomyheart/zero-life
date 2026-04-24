package request

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
}

type RecurringTransactionListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Active   *bool  `form:"active"`
}