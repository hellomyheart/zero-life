package model

import (
	"time"

	"gorm.io/gorm"
)

type LogicType string

const (
	LogicTypeAnd LogicType = "and"
	LogicTypeOr  LogicType = "or"
)

type RuleTrigger string

const (
	RuleTriggerOnCreate RuleTrigger = "on_create"
	RuleTriggerOnUpdate RuleTrigger = "on_update"
)

type Rule struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Priority  int            `gorm:"default:0" json:"priority"`
	IsEnabled bool           `gorm:"default:true" json:"is_enabled"`
	LogicType LogicType      `gorm:"not null;size:20" json:"logic_type"`
	Trigger   RuleTrigger    `gorm:"not null;size:20;default:on_create" json:"trigger"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Conditions []RuleCondition `gorm:"foreignKey:RuleID" json:"conditions,omitempty"`
	Actions    []RuleAction    `gorm:"foreignKey:RuleID" json:"actions,omitempty"`
}

func (Rule) TableName() string { return "rules" }

type ConditionField string

const (
	ConditionFieldDescription        ConditionField = "description"
	ConditionFieldAmount             ConditionField = "amount"
	ConditionFieldSourceAccount      ConditionField = "source_account"
	ConditionFieldDestinationAccount ConditionField = "destination_account"
	ConditionFieldCategory           ConditionField = "category"
	ConditionFieldTag                ConditionField = "tag"
	ConditionFieldTransactionType    ConditionField = "transaction_type"
	ConditionFieldBudget             ConditionField = "budget"
	ConditionFieldBill               ConditionField = "bill"
	ConditionFieldNotes              ConditionField = "notes"
	ConditionFieldDateAfter          ConditionField = "date_after"
	ConditionFieldDateBefore         ConditionField = "date_before"
)

type ConditionOperator string

const (
	OperatorContains    ConditionOperator = "contains"
	OperatorEquals      ConditionOperator = "equals"
	OperatorStartsWith  ConditionOperator = "starts_with"
	OperatorEndsWith    ConditionOperator = "ends_with"
	OperatorNotContains ConditionOperator = "not_contains"
	OperatorNotEquals   ConditionOperator = "not_equals"
	OperatorLess        ConditionOperator = "less"
	OperatorMore        ConditionOperator = "more"
	OperatorIsEmpty     ConditionOperator = "is_empty"
	OperatorIsNotEmpty  ConditionOperator = "is_not_empty"
)

type RuleCondition struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleID    uint64            `gorm:"not null;index" json:"rule_id"`
	Field     ConditionField    `gorm:"not null;size:20" json:"field"`
	Operator  ConditionOperator `gorm:"not null;size:20" json:"operator"`
	Value     string            `gorm:"not null;size:255" json:"value"`
	CreatedAt time.Time         `gorm:"not null" json:"created_at"`
}

func (RuleCondition) TableName() string { return "rule_conditions" }

type ActionType string

const (
	ActionTypeSetCategory   ActionType = "set_category"
	ActionTypeAddTag        ActionType = "add_tag"
	ActionTypeSetNotes      ActionType = "set_notes"
	ActionTypeSetBudget     ActionType = "set_budget"
	ActionTypeRemoveTag     ActionType = "remove_tag"
	ActionTypeSetDescription ActionType = "set_description"
	ActionTypeClearCategory ActionType = "clear_category"
	ActionTypeClearBudget   ActionType = "clear_budget"
	ActionTypeClearNotes    ActionType = "clear_notes"
	ActionTypeAppendNotes   ActionType = "append_notes"
	ActionTypePrependNotes  ActionType = "prepend_notes"
)

type RuleAction struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleID    uint64     `gorm:"not null;index" json:"rule_id"`
	Type      ActionType `gorm:"not null;size:20" json:"type"`
	Value     string     `gorm:"not null;size:255" json:"value"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
}

func (RuleAction) TableName() string { return "rule_actions" }
