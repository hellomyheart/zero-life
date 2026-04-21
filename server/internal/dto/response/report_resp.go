package response

type IncomeExpenseResp struct {
	TotalIncome  string              `json:"total_income"`
	TotalExpense string              `json:"total_expense"`
	NetIncome    string              `json:"net_income"`
	PrevIncome   string              `json:"prev_income,omitempty"`
	PrevExpense  string              `json:"prev_expense,omitempty"`
	PrevNet      string              `json:"prev_net,omitempty"`
	Details      []PeriodDetailResp  `json:"details,omitempty"`
}

type PeriodDetailResp struct {
	Period  string `json:"period"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
	Net     string `json:"net"`
}

type CategoryReportResp struct {
	ExpenseByCategory []CategoryItemResp `json:"expense_by_category"`
	IncomeByCategory  []CategoryItemResp `json:"income_by_category"`
}

type CategoryItemResp struct {
	CategoryID   uint64 `json:"category_id"`
	CategoryName string `json:"category_name"`
	Amount       string `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type BudgetReportResp struct {
	Items []BudgetReportItemResp `json:"items"`
}

type BudgetReportItemResp struct {
	BudgetID   uint64 `json:"budget_id"`
	BudgetName string `json:"budget_name"`
	Amount     string `json:"amount"`
	Spent      string `json:"spent"`
	Remaining  string `json:"remaining"`
	UsageRate  float64 `json:"usage_rate"`
}

type NetWorthResp struct {
	TotalNetWorth string              `json:"total_net_worth"`
	Trend         []NetWorthPointResp `json:"trend"`
}

type NetWorthPointResp struct {
	Date      string `json:"date"`
	NetWorth  string `json:"net_worth"`
}

type TrendResp struct {
	Items []TrendItemResp `json:"items"`
}

type TrendItemResp struct {
	Period  string `json:"period"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
}
