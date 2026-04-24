// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// LogicType 条件逻辑类型
type LogicType string

const (
	LogicTypeAnd LogicType = "and" // 所有条件都满足才匹配
	LogicTypeOr  LogicType = "or"  // 任一条件满足即匹配
)

// RuleTrigger 规则触发时机
type RuleTrigger string

const (
	RuleTriggerOnCreate RuleTrigger = "on_create" // 创建交易时触发
	RuleTriggerOnUpdate RuleTrigger = "on_update" // 更新交易时触发
)

// Rule 规则模型，对应rules表
// 定义自动化规则，当交易满足条件时自动执行操作
type Rule struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`           // 所属用户ID
	GroupID   uint64         `gorm:"index" json:"group_id"`                   // 所属规则组ID
	Name      string         `gorm:"not null;size:100" json:"name"`           // 规则名称
	Priority  int            `gorm:"default:0" json:"priority"`               // 优先级，数值越小越先执行
	IsEnabled bool           `gorm:"default:true" json:"is_enabled"`          // 是否启用
	LogicType LogicType      `gorm:"not null;size:20" json:"logic_type"`      // 条件逻辑类型（and/or）
	Trigger   RuleTrigger    `gorm:"not null;size:20;default:on_create" json:"trigger"` // 触发时机
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间

	Group      RuleGroup      `gorm:"foreignKey:GroupID" json:"group,omitempty"`       // 所属规则组
	Conditions []RuleCondition `gorm:"foreignKey:RuleID" json:"conditions,omitempty"`   // 条件列表
	Actions    []RuleAction    `gorm:"foreignKey:RuleID" json:"actions,omitempty"`       // 操作列表
}

// TableName 指定表名
func (Rule) TableName() string { return "rules" }

// ConditionField 条件字段枚举
type ConditionField string

const (
	ConditionFieldDescription        ConditionField = "description"         // 交易描述
	ConditionFieldAmount             ConditionField = "amount"              // 交易金额
	ConditionFieldSourceAccount      ConditionField = "source_account"      // 源账户
	ConditionFieldDestinationAccount ConditionField = "destination_account" // 目标账户
	ConditionFieldCategory           ConditionField = "category"           // 分类
	ConditionFieldTag                ConditionField = "tag"                // 标签
	ConditionFieldTransactionType    ConditionField = "transaction_type"   // 交易类型
	ConditionFieldBudget             ConditionField = "budget"             // 预算
	ConditionFieldBill               ConditionField = "bill"               // 账单
	ConditionFieldNotes              ConditionField = "notes"              // 备注
	ConditionFieldDateAfter          ConditionField = "date_after"         // 日期晚于
	ConditionFieldDateBefore         ConditionField = "date_before"        // 日期早于
)

// ConditionOperator 条件运算符枚举
type ConditionOperator string

const (
	OperatorContains    ConditionOperator = "contains"     // 包含
	OperatorEquals      ConditionOperator = "equals"       // 等于
	OperatorStartsWith  ConditionOperator = "starts_with"  // 以...开头
	OperatorEndsWith    ConditionOperator = "ends_with"    // 以...结尾
	OperatorNotContains ConditionOperator = "not_contains"  // 不包含
	OperatorNotEquals   ConditionOperator = "not_equals"    // 不等于
	OperatorLess        ConditionOperator = "less"          // 小于
	OperatorMore        ConditionOperator = "more"          // 大于
	OperatorIsEmpty     ConditionOperator = "is_empty"     // 为空
	OperatorIsNotEmpty  ConditionOperator = "is_not_empty"  // 不为空
)

// RuleCondition 规则条件模型，对应rule_conditions表
type RuleCondition struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleID    uint64            `gorm:"not null;index" json:"rule_id"`               // 所属规则ID
	Field     ConditionField    `gorm:"not null;size:20" json:"field"`               // 条件字段
	Operator  ConditionOperator `gorm:"not null;size:20" json:"operator"`            // 条件运算符
	Value     string            `gorm:"not null;size:255" json:"value"`              // 条件值
	CreatedAt time.Time         `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (RuleCondition) TableName() string { return "rule_conditions" }

// ActionType 操作类型枚举
type ActionType string

const (
	ActionTypeSetCategory    ActionType = "set_category"     // 设置分类
	ActionTypeAddTag         ActionType = "add_tag"          // 添加标签
	ActionTypeSetNotes       ActionType = "set_notes"        // 设置备注
	ActionTypeSetBudget      ActionType = "set_budget"       // 设置预算
	ActionTypeRemoveTag      ActionType = "remove_tag"       // 移除标签
	ActionTypeSetDescription ActionType = "set_description"  // 设置描述
	ActionTypeClearCategory  ActionType = "clear_category"   // 清除分类
	ActionTypeClearBudget    ActionType = "clear_budget"     // 清除预算
	ActionTypeClearNotes     ActionType = "clear_notes"      // 清除备注
	ActionTypeAppendNotes    ActionType = "append_notes"     // 追加备注
	ActionTypePrependNotes   ActionType = "prepend_notes"    // 前置备注
)

// RuleAction 规则操作模型，对应rule_actions表
type RuleAction struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleID    uint64     `gorm:"not null;index" json:"rule_id"`      // 所属规则ID
	Type      ActionType `gorm:"not null;size:20" json:"type"`       // 操作类型
	Value     string     `gorm:"not null;size:255" json:"value"`     // 操作值
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (RuleAction) TableName() string { return "rule_actions" }
