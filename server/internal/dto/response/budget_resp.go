// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// BudgetResp 预算响应
// 包含预算基本信息和执行情况（已支出、剩余、使用率、状态）
type BudgetResp struct {
	ID         uint64         `json:"id"`          // 预算ID
	Name       string         `json:"name"`        // 预算名称
	Amount     string         `json:"amount"`      // 预算金额
	Period     string         `json:"period"`      // 预算周期
	IsEnabled  bool           `json:"is_enabled"`  // 是否启用
	Categories []CategoryResp `json:"categories"`  // 关联分类列表
	Spent      string         `json:"spent"`       // 已支出金额
	Remaining  string         `json:"remaining"`   // 剩余金额
	UsageRate  float64        `json:"usage_rate"`  // 使用率（比率，0~1+，如0.8表示80%）
	Status     string         `json:"status"`      // 状态：normal(正常)/warning(预警)/overspent(超支)
	CreatedAt  time.Time      `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"`  // 更新时间
}

// BudgetHistoryResp 预算历史记录响应
type BudgetHistoryResp struct {
	ID          uint64    `json:"id"`           // 记录ID
	PeriodStart time.Time `json:"period_start"` // 周期开始时间
	PeriodEnd   time.Time `json:"period_end"`   // 周期结束时间
	Amount      string    `json:"amount"`       // 预算金额
	Spent       string    `json:"spent"`        // 实际支出
	UsageRate   float64   `json:"usage_rate"`   // 使用率（比率，0~1+，如0.8表示80%）
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}
