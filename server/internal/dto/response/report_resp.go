// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// IncomeExpenseResp 收支报表响应
type IncomeExpenseResp struct {
	TotalIncome  string             `json:"total_income"`            // 总收入
	TotalExpense string             `json:"total_expense"`           // 总支出
	NetIncome    string             `json:"net_income"`              // 净收入
	PrevIncome   string             `json:"prev_income,omitempty"`   // 上期收入（用于对比）
	PrevExpense  string             `json:"prev_expense,omitempty"`  // 上期支出
	PrevNet      string             `json:"prev_net,omitempty"`      // 上期净收入
	Details      []PeriodDetailResp `json:"details,omitempty"`       // 分时段明细
}

// PeriodDetailResp 时段明细响应
type PeriodDetailResp struct {
	Period  string `json:"period"`  // 时段标识
	Income  string `json:"income"`  // 收入
	Expense string `json:"expense"` // 支出
	Net     string `json:"net"`     // 净值
}

// CategoryReportResp 分类报表响应
type CategoryReportResp struct {
	ExpenseByCategory []CategoryItemResp `json:"expense_by_category"` // 按分类的支出分布
	IncomeByCategory  []CategoryItemResp `json:"income_by_category"`  // 按分类的收入分布
}

// CategoryItemResp 分类报表项响应
type CategoryItemResp struct {
	CategoryID   uint64  `json:"category_id"`   // 分类ID
	CategoryName string  `json:"category_name"` // 分类名称
	Amount       string  `json:"amount"`        // 金额
	Percentage   float64 `json:"percentage"`    // 占比百分比
}

// BudgetReportResp 预算报表响应
type BudgetReportResp struct {
	Items []BudgetReportItemResp `json:"items"` // 预算报表项列表
}

// BudgetReportItemResp 预算报表项响应
type BudgetReportItemResp struct {
	BudgetID   uint64  `json:"budget_id"`   // 预算ID
	BudgetName string  `json:"budget_name"` // 预算名称
	Amount     string  `json:"amount"`      // 预算金额
	Spent      string  `json:"spent"`       // 已支出金额
	Remaining  string  `json:"remaining"`   // 剩余金额
	UsageRate  float64 `json:"usage_rate"`  // 使用率
}

// NetWorthResp 净资产报表响应
type NetWorthResp struct {
	TotalNetWorth string              `json:"total_net_worth"` // 总净资产
	Trend         []NetWorthPointResp `json:"trend"`           // 净资产趋势
}

// NetWorthPointResp 净资产趋势点响应
type NetWorthPointResp struct {
	Date     string `json:"date"`      // 日期
	NetWorth string `json:"net_worth"` // 净资产值
}

// TrendResp 趋势报表响应
type TrendResp struct {
	Items []TrendItemResp `json:"items"` // 趋势数据列表
}

// TrendItemResp 趋势报表项响应
type TrendItemResp struct {
	Period  string `json:"period"`  // 时段标识
	Income  string `json:"income"`  // 收入
	Expense string `json:"expense"` // 支出
}

// TagReportResp 标签报表响应
type TagReportResp struct {
	Items []TagReportItemResp `json:"items"` // 标签报表项列表
}

// TagReportItemResp 标签报表项响应
type TagReportItemResp struct {
	TagID   uint64 `json:"tag_id"`   // 标签ID
	TagName string `json:"tag_name"` // 标签名称
	Income  string `json:"income"`   // 收入
	Expense string `json:"expense"`  // 支出
}

// AuditReportResp 审计报表响应
type AuditReportResp struct {
	AccountID      uint64                `json:"account_id"`       // 账户ID
	AccountName    string                `json:"account_name"`     // 账户名称
	InitialBalance string                `json:"initial_balance"`  // 期初余额
	FinalBalance   string                `json:"final_balance"`    // 期末余额
	StartDate      string                `json:"start_date"`       // 开始日期
	EndDate        string                `json:"end_date"`         // 结束日期
	Items          []AuditReportItemResp `json:"items"`            // 审计明细列表
}

// AuditReportItemResp 审计报表项响应
type AuditReportItemResp struct {
	TransactionID  uint64    `json:"transaction_id"`            // 交易ID
	Date           time.Time `json:"date"`                     // 交易日期
	Description    string    `json:"description"`              // 交易描述
	Type           string    `json:"type"`                     // 交易类型
	Amount         string    `json:"amount"`                   // 交易金额
	RunningBalance string    `json:"running_balance"`          // 余额流水
	CategoryName   string    `json:"category_name,omitempty"`  // 分类名称
	IsReconciled   bool      `json:"is_reconciled"`            // 是否已对账
}
