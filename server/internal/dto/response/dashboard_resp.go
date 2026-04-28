// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// DashboardResp 仪表盘响应
// 包含首页仪表盘的所有汇总数据
type DashboardResp struct {
	MonthIncome    string             `json:"month_income"`     // 本月收入
	MonthExpense   string             `json:"month_expense"`    // 本月支出
	NetIncome      string             `json:"net_income"`       // 净收入
	TotalBalance   string             `json:"total_balance"`    // 总资产余额
	BudgetAlerts   []BudgetAlertResp  `json:"budget_alerts"`    // 预算预警列表
	RecurringReminders  []RecurringReminderResp `json:"recurring_reminders"`   // 循环交易提醒列表
	RecentTxns     []TransactionResp  `json:"recent_txns"`      // 最近交易列表
}

// BudgetAlertResp 预算预警响应
type BudgetAlertResp struct {
	BudgetID   uint64  `json:"budget_id"`   // 预算ID
	BudgetName string  `json:"budget_name"` // 预算名称
	Spent      string  `json:"spent"`       // 已支出金额
	Amount     string  `json:"amount"`      // 预算金额
	UsageRate  float64 `json:"usage_rate"`  // 使用率
	Status     string  `json:"status"`      // 状态
}

// RecurringReminderResp 循环交易提醒响应
type RecurringReminderResp struct {
	RecurringID uint64 `json:"recurring_id"` // 循环交易ID
	Name        string `json:"name"`         // 名称
	Amount      string `json:"amount"`       // 金额
	NextDue     string `json:"next_due"`     // 下次到期日
}
