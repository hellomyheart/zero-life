// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// RecurringTransactionResp 循环交易响应
type RecurringTransactionResp struct {
	ID             uint64     `json:"id"`                          // 循环交易ID
	Description    string     `json:"description"`                 // 交易描述
	Amount         string     `json:"amount"`                      // 交易金额
	SourceID       uint64     `json:"source_id"`                   // 源账户ID
	DestinationID  *uint64    `json:"destination_id"`              // 目标账户ID
	CategoryID     *uint64    `json:"category_id"`                 // 分类ID
	Notes          string     `json:"notes"`                       // 备注
	RecurrenceType string     `json:"recurrence_type"`             // 重复类型
	RepeatEvery    int        `json:"repeat_every"`                // 重复间隔
	StartDate      time.Time  `json:"start_date"`                  // 开始日期
	EndDate        *time.Time `json:"end_date,omitempty"`          // 结束日期
	NextOccurrence time.Time  `json:"next_occurrence"`             // 下次执行日期
	IsActive       bool       `json:"is_active"`                   // 是否启用
	CreatedAt      time.Time  `json:"created_at"`                  // 创建时间
	UpdatedAt      time.Time  `json:"updated_at"`                  // 更新时间
}

// RecurringTransactionLogResp 循环交易执行日志响应
type RecurringTransactionLogResp struct {
	ID                  uint64    `json:"id"`                     // 日志ID
	RecurringTransactionID uint64 `json:"recurring_transaction_id"` // 循环交易ID
	TransactionID       uint64    `json:"transaction_id"`          // 创建的交易ID
	OccurrenceDate      time.Time `json:"occurrence_date"`         // 执行日期
	CreatedAt           time.Time `json:"created_at"`              // 创建时间
}