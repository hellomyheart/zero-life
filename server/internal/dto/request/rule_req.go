// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateRuleReq 创建规则请求
// 规则用于在交易创建或更新时自动执行操作
type CreateRuleReq struct {
	Name       string             `json:"name" binding:"required"`                          // 规则名称
	Priority   int                `json:"priority"`                                         // 优先级，数值越小越先执行
	IsEnabled  bool               `json:"is_enabled"`                                       // 是否启用
	LogicType  string             `json:"logic_type" binding:"required,oneof=and or"`       // 条件逻辑：and(全部满足)/or(任一满足)
	Trigger    string             `json:"trigger" binding:"required,oneof=on_create on_update"` // 触发时机：on_create(创建时)/on_update(更新时)
	Conditions []RuleConditionReq `json:"conditions" binding:"required,min=1"`              // 条件列表，至少1个
	Actions    []RuleActionReq    `json:"actions" binding:"required,min=1"`                 // 动作列表，至少1个
}

// UpdateRuleReq 更新规则请求
type UpdateRuleReq struct {
	Name       string             `json:"name"`
	Priority   *int               `json:"priority"`
	IsEnabled  *bool              `json:"is_enabled"`
	LogicType  string             `json:"logic_type" binding:"omitempty,oneof=and or"`
	Trigger    string             `json:"trigger" binding:"omitempty,oneof=on_create on_update"`
	Conditions []RuleConditionReq `json:"conditions"`
	Actions    []RuleActionReq    `json:"actions"`
}

// RuleConditionReq 规则条件请求
// 定义规则匹配的条件
type RuleConditionReq struct {
	Field    string `json:"field" binding:"required,oneof=description amount source_account destination_account category tag transaction_type budget bill notes date_after date_before"` // 匹配字段
	Operator string `json:"operator" binding:"required,oneof=contains equals starts_with ends_with not_contains not_equals less more is_empty is_not_empty"` // 匹配运算符
	Value    string `json:"value"` // 匹配值
}

// RuleActionReq 规则动作请求
// 定义规则匹配后执行的操作
type RuleActionReq struct {
	Type  string `json:"type" binding:"required,oneof=set_category add_tag set_notes set_budget remove_tag set_description clear_category clear_budget clear_notes append_notes prepend_notes"` // 动作类型
	Value string `json:"value"` // 动作值
}

// ExecuteRuleReq 执行规则请求
type ExecuteRuleReq struct {
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}
