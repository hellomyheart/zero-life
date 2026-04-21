package request

type CreateRuleReq struct {
	Name       string              `json:"name" binding:"required"`
	Priority   int                 `json:"priority"`
	IsEnabled  bool                `json:"is_enabled"`
	LogicType  string              `json:"logic_type" binding:"required,oneof=and or"`
	Conditions []RuleConditionReq  `json:"conditions" binding:"required,min=1"`
	Actions    []RuleActionReq     `json:"actions" binding:"required,min=1"`
}

type UpdateRuleReq struct {
	Name       string              `json:"name"`
	Priority   *int                `json:"priority"`
	IsEnabled  *bool               `json:"is_enabled"`
	LogicType  string              `json:"logic_type" binding:"omitempty,oneof=and or"`
	Conditions []RuleConditionReq  `json:"conditions"`
	Actions    []RuleActionReq     `json:"actions"`
}

type RuleConditionReq struct {
	Field    string `json:"field" binding:"required,oneof=description amount source_account"`
	Operator string `json:"operator" binding:"required,oneof=contains equals starts_with ends_with"`
	Value    string `json:"value" binding:"required"`
}

type RuleActionReq struct {
	Type  string `json:"type" binding:"required,oneof=set_category add_tag set_notes"`
	Value string `json:"value" binding:"required"`
}

type ExecuteRuleReq struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
