// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// BillResp 账单响应
type BillResp struct {
	ID         uint64    `json:"id"`          // 账单ID
	Name       string    `json:"name"`        // 账单名称
	Amount     string    `json:"amount"`      // 账单金额
	RepeatRule string    `json:"repeat_rule"` // 重复规则
	NextDue    time.Time `json:"next_due"`    // 下次到期日
	SourceID   *uint64   `json:"source_id"`   // 支出账户ID
	CategoryID *uint64   `json:"category_id"` // 分类ID
	Notes      string    `json:"notes"`       // 备注
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}
