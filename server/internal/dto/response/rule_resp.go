// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// RuleResp 规则响应
type RuleResp struct {
	ID         uint64              `json:"id"`         // 规则ID
	Name       string              `json:"name"`       // 规则名称
	Priority   int                 `json:"priority"`   // 优先级
	IsEnabled  bool                `json:"is_enabled"` // 是否启用
	LogicType  string              `json:"logic_type"` // 条件逻辑类型
	Trigger    string              `json:"trigger"`    // 触发时机
	Conditions []RuleConditionResp `json:"conditions"` // 条件列表
	Actions    []RuleActionResp    `json:"actions"`    // 动作列表
	CreatedAt  time.Time           `json:"created_at"` // 创建时间
	UpdatedAt  time.Time           `json:"updated_at"` // 更新时间
}

// RuleConditionResp 规则条件响应
type RuleConditionResp struct {
	ID       uint64 `json:"id"`       // 条件ID
	Field    string `json:"field"`    // 匹配字段
	Operator string `json:"operator"` // 匹配运算符
	Value    string `json:"value"`    // 匹配值
}

// RuleActionResp 规则动作响应
type RuleActionResp struct {
	ID    uint64 `json:"id"`    // 动作ID
	Type  string `json:"type"`  // 动作类型
	Value string `json:"value"` // 动作值
}

// RuleExecuteResultResp 规则执行结果响应
type RuleExecuteResultResp struct {
	MatchedCount int      `json:"matched_count"` // 匹配的交易数
	SuccessCount int      `json:"success_count"` // 成功执行的动作数
	FailCount    int      `json:"fail_count"`    // 失败的动作数
	Errors       []string `json:"errors,omitempty"` // 错误信息列表
}
