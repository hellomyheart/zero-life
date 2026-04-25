// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// ReconciliationResp 对账响应
type ReconciliationResp struct {
	ID              uint64    `json:"id"`               // 对账ID
	AccountID       uint64    `json:"account_id"`       // 对账账户ID
	StartDate       time.Time `json:"start_date"`       // 对账开始日期
	EndDate         time.Time `json:"end_date"`         // 对账结束日期
	StartingBalance string    `json:"starting_balance"` // 期初余额
	EndingBalance   string    `json:"ending_balance"`   // 期末余额
	BookBalance     string    `json:"book_balance"`     // 账面余额
	Difference      string    `json:"difference"`       // 差额（期末余额 - 账面余额）
	Status          string    `json:"status"`           // 状态：open/closed
	CreatedAt       time.Time `json:"created_at"`       // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`       // 更新时间
}

// ReconciliationEntryResp 对账条目响应
type ReconciliationEntryResp struct {
	ID               uint64 `json:"id"`                 // 条目ID
	ReconciliationID uint64 `json:"reconciliation_id"`  // 对账ID
	TransactionID    uint64 `json:"transaction_id"`     // 交易ID
	Matched          bool   `json:"matched"`            // 是否已匹配
	AmountDifference string `json:"amount_difference"`  // 金额差异
}