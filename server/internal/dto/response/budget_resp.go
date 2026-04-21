package response

import "time"

type BudgetResp struct {
	ID         uint64         `json:"id"`
	Name       string         `json:"name"`
	Amount     string         `json:"amount"`
	Period     string         `json:"period"`
	IsEnabled  bool           `json:"is_enabled"`
	Categories []CategoryResp `json:"categories"`
	Spent      string         `json:"spent"`
	Remaining  string         `json:"remaining"`
	UsageRate  float64        `json:"usage_rate"`
	Status     string         `json:"status"` // normal, warning, overspent
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type BudgetHistoryResp struct {
	ID          uint64    `json:"id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Amount      string    `json:"amount"`
	Spent       string    `json:"spent"`
	UsageRate   float64   `json:"usage_rate"`
	CreatedAt   time.Time `json:"created_at"`
}
