package response

type DashboardResp struct {
	MonthIncome    string              `json:"month_income"`
	MonthExpense   string              `json:"month_expense"`
	NetIncome      string              `json:"net_income"`
	TotalBalance   string              `json:"total_balance"`
	BudgetAlerts   []BudgetAlertResp   `json:"budget_alerts"`
	BillReminders  []BillReminderResp  `json:"bill_reminders"`
	RecentTxns     []TransactionResp    `json:"recent_txns"`
}

type BudgetAlertResp struct {
	BudgetID   uint64 `json:"budget_id"`
	BudgetName string `json:"budget_name"`
	Spent      string `json:"spent"`
	Amount     string `json:"amount"`
	UsageRate  float64 `json:"usage_rate"`
	Status     string  `json:"status"`
}

type BillReminderResp struct {
	BillID   uint64 `json:"bill_id"`
	BillName string `json:"bill_name"`
	Amount   string `json:"amount"`
	NextDue  string `json:"next_due"`
}
