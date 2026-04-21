package response

import "time"

type RuleResp struct {
	ID         uint64              `json:"id"`
	Name       string              `json:"name"`
	Priority   int                 `json:"priority"`
	IsEnabled  bool                `json:"is_enabled"`
	LogicType  string              `json:"logic_type"`
	Conditions []RuleConditionResp `json:"conditions"`
	Actions    []RuleActionResp    `json:"actions"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type RuleConditionResp struct {
	ID       uint64 `json:"id"`
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type RuleActionResp struct {
	ID    uint64 `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
