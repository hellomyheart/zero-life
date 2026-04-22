package request

type CreateRuleReq struct {
	Name       string             `json:"name" binding:"required"`
	Priority   int                `json:"priority"`
	IsEnabled  bool               `json:"is_enabled"`
	LogicType  string             `json:"logic_type" binding:"required,oneof=and or"`
	Trigger    string             `json:"trigger" binding:"required,oneof=on_create on_update"`
	Conditions []RuleConditionReq `json:"conditions" binding:"required,min=1"`
	Actions    []RuleActionReq    `json:"actions" binding:"required,min=1"`
}

type UpdateRuleReq struct {
	Name       string             `json:"name"`
	Priority   *int               `json:"priority"`
	IsEnabled  *bool              `json:"is_enabled"`
	LogicType  string             `json:"logic_type" binding:"omitempty,oneof=and or"`
	Trigger    string             `json:"trigger" binding:"omitempty,oneof=on_create on_update"`
	Conditions []RuleConditionReq `json:"conditions"`
	Actions    []RuleActionReq    `json:"actions"`
}

type RuleConditionReq struct {
	Field    string `json:"field" binding:"required,oneof=description amount source_account destination_account category tag transaction_type budget bill notes date_after date_before"`
	Operator string `json:"operator" binding:"required,oneof=contains equals starts_with ends_with not_contains not_equals less more is_empty is_not_empty"`
	Value    string `json:"value"`
}

type RuleActionReq struct {
	Type  string `json:"type" binding:"required,oneof=set_category add_tag set_notes set_budget remove_tag set_description clear_category clear_budget clear_notes append_notes prepend_notes"`
	Value string `json:"value"`
}

type ExecuteRuleReq struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
