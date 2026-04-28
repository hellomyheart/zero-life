// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// RecurringTransactionResp 循环交易响应
type RecurringTransactionResp struct {
	ID             uint64     `json:"id"`
	Description    string     `json:"description"`
	Amount         string     `json:"amount"`
	SourceID       uint64     `json:"source_id"`
	DestinationID  *uint64    `json:"destination_id"`
	CategoryID     *uint64    `json:"category_id"`
	Notes          string     `json:"notes"`
	RecurrenceType string     `json:"recurrence_type"`
	RepeatEvery    int        `json:"repeat_every"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	NextOccurrence time.Time  `json:"next_occurrence"`
	IsActive       bool       `json:"is_active"`
	ReminderDays   int        `json:"reminder_days"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// RecurringTransactionLogResp 循环交易执行日志响应
type RecurringTransactionLogResp struct {
	ID                  uint64    `json:"id"`                     // 日志ID
	RecurringTransactionID uint64 `json:"recurring_transaction_id"` // 循环交易ID
	TransactionID       uint64    `json:"transaction_id"`          // 创建的交易ID
	OccurrenceDate      time.Time `json:"occurrence_date"`         // 执行日期
	CreatedAt           time.Time `json:"created_at"`              // 创建时间
}