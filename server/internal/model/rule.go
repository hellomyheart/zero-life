// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单、规则等
package model

import (
	"time"

	"gorm.io/gorm"
)

// LogicType 条件逻辑类型
// 定义多个条件之间的逻辑关系
type LogicType string

const (
	// LogicTypeAnd 所有条件都满足才匹配
	// 即"且"关系，条件必须全部成立
	// 如：描述包含"星巴克" AND 金额小于 100
	LogicTypeAnd LogicType = "and"

	// LogicTypeOr 任一条件满足即匹配
	// 即"或"关系，条件只需一个成立
	// 如：描述包含"星巴克" OR 描述包含"Costa"
	LogicTypeOr LogicType = "or"
)

// RuleTrigger 规则触发时机
// 定义规则在什么时候执行
type RuleTrigger string

const (
	// RuleTriggerOnCreate 创建交易时触发
	// 当新建一笔交易时，自动检查并执行规则
	// 如：创建交易时自动设置分类
	RuleTriggerOnCreate RuleTrigger = "on_create"

	// RuleTriggerOnUpdate 更新交易时触发
	// 当修改一笔交易时，自动检查并执行规则
	// 如：更新交易时自动打标签
	RuleTriggerOnUpdate RuleTrigger = "on_update"
)

// Rule 规则模型，对应 rules 表
//
// 功能说明：
// - 定义自动化规则，当交易满足条件时自动执行操作
// - 支持多条件组合（AND/OR 逻辑）
// - 支持多种触发时机（创建/更新交易时）
// - 支持优先级排序，数值越小越先执行
// - 可归入规则组进行分组管理
//
// 使用场景：
// - 自动分类：描述包含"星巴克" -> 设置分类为"餐饮/咖啡"
// - 自动打标签：金额大于 1000 -> 添加标签"大额支出"
// - 自动设置备注：描述包含"滴滴" -> 设置备注"打车"
//
// 规则执行流程：
// 1. 交易创建或更新时，触发规则引擎
// 2. 按 Priority 排序，依次检查每条规则
// 3. 根据LogicType（AND/OR）判断条件是否满足
// 4. 条件满足则执行所有关联的操作
// 5. 继续检查下一条规则（除非操作中指定停止）
//
// 与其他模型的关系：
// - RuleGroup: 多对一关系，规则属于一个规则组
// - RuleCondition: 一对多关系，规则包含多个条件
// - RuleAction: 一对多关系，规则包含多个操作
//
// 示例：
//   规则：星巴克自动分类
//   条件：描述包含"星巴克"（AND 逻辑）
//   操作：设置分类为"餐饮/咖啡"，添加标签"咖啡"
//   触发时机：创建交易时
//   优先级：10
type Rule struct {
	// ID 规则唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的规则集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询规则
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// GroupID 所属规则组 ID
	// 关联到 rule_groups 表，用于分组管理规则
	// gorm:"index" 创建索引，加速按组查询规则
	GroupID uint64 `gorm:"index" json:"group_id"`

	// Name 规则名称
	// 用户自定义的规则名称，便于识别
	// 如："星巴克自动分类"、"大额支出标记"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// Priority 优先级
	// 数值越小越先执行
	// 如：Priority=1 的规则先于 Priority=10 的规则执行
	// 默认值：0
	Priority int `gorm:"default:0" json:"priority"`

	// IsEnabled 是否启用
	// true: 规则生效，满足条件时执行操作
	// false: 规则暂停，不执行
	// 默认值：true
	IsEnabled bool `gorm:"default:true" json:"is_enabled"`

	// LogicType 条件逻辑类型
	// 取值为 LogicType 枚举：and/or
	// and: 所有条件都满足才执行操作
	// or: 任一条件满足即执行操作
	// gorm:"size:20" 限制最大长度为 20 个字符
	LogicType LogicType `gorm:"not null;size:20" json:"logic_type"`

	// Trigger 触发时机
	// 取值为 RuleTrigger 枚举：on_create/on_update
	// on_create: 创建交易时触发
	// on_update: 更新交易时触发
	// gorm:"size:20" 限制最大长度为 20 个字符
	// 默认值：on_create
	Trigger RuleTrigger `gorm:"not null;size:20;default:on_create" json:"trigger"`

	// CreatedAt 规则创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 规则最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Group 所属规则组
	// 通过 GroupID 外键关联到 rule_groups 表
	// json:"group,omitempty" 当字段为零值时 JSON 序列化时忽略
	Group RuleGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`

	// Conditions 条件列表
	// 通过 RuleID 外键关联到 rule_conditions 表
	// 一条规则可以有多个条件，条件之间通过 LogicType 组合
	// json:"conditions,omitempty" 当字段为零值时 JSON 序列化时忽略
	Conditions []RuleCondition `gorm:"foreignKey:RuleID" json:"conditions,omitempty"`

	// Actions 操作列表
	// 通过 RuleID 外键关联到 rule_actions 表
	// 条件满足时，依次执行所有关联的操作
	// json:"actions,omitempty" 当字段为零值时 JSON 序列化时忽略
	Actions []RuleAction `gorm:"foreignKey:RuleID" json:"actions,omitempty"`
}

// TableName 指定 Rule 模型对应的数据库表名为 rules
func (Rule) TableName() string { return "rules" }

// ConditionField 条件字段枚举
// 定义规则条件可以匹配的交易字段
type ConditionField string

const (
	// ConditionFieldDescription 匹配交易描述
	// 如：描述包含"星巴克"
	ConditionFieldDescription ConditionField = "description"

	// ConditionFieldAmount 匹配交易金额
	// 如：金额大于 100
	ConditionFieldAmount ConditionField = "amount"

	// ConditionFieldSourceAccount 匹配源账户
	// 如：源账户为"招商银行卡"
	ConditionFieldSourceAccount ConditionField = "source_account"

	// ConditionFieldDestinationAccount 匹配目标账户
	// 如：目标账户为"餐饮支出"
	ConditionFieldDestinationAccount ConditionField = "destination_account"

	// ConditionFieldCategory 匹配分类
	// 如：分类为"餐饮"
	ConditionFieldCategory ConditionField = "category"

	// ConditionFieldTag 匹配标签
	// 如：标签包含"出差"
	ConditionFieldTag ConditionField = "tag"

	// ConditionFieldTransactionType 匹配交易类型
	// 如：交易类型为"支出"
	ConditionFieldTransactionType ConditionField = "transaction_type"

	// ConditionFieldBudget 匹配预算
	// 如：关联预算为"月度餐饮预算"
	ConditionFieldBudget ConditionField = "budget"

	// ConditionFieldBill 匹配循环交易
	// 如：关联循环交易为"房租"
	ConditionFieldBill ConditionField = "bill"

	// ConditionFieldNotes 匹配备注
	// 如：备注包含"报销"
	ConditionFieldNotes ConditionField = "notes"

	// ConditionFieldDateAfter 匹配日期晚于
	// 如：交易日期晚于 2026-04-01
	ConditionFieldDateAfter ConditionField = "date_after"

	// ConditionFieldDateBefore 匹配日期早于
	// 如：交易日期早于 2026-04-30
	ConditionFieldDateBefore ConditionField = "date_before"
)

// ConditionOperator 条件运算符枚举
// 定义条件字段与条件值之间的比较方式
type ConditionOperator string

const (
	// OperatorContains 包含：字段值包含指定字符串
	// 如：描述 contains "星巴克"
	OperatorContains ConditionOperator = "contains"

	// OperatorEquals 等于：字段值完全等于指定值
	// 如：金额 equals "100"
	OperatorEquals ConditionOperator = "equals"

	// OperatorStartsWith 以...开头：字段值以指定字符串开头
	// 如：描述 starts_with "滴滴"
	OperatorStartsWith ConditionOperator = "starts_with"

	// OperatorEndsWith 以...结尾：字段值以指定字符串结尾
	// 如：描述 ends_with "退款"
	OperatorEndsWith ConditionOperator = "ends_with"

	// OperatorNotContains 不包含：字段值不包含指定字符串
	OperatorNotContains ConditionOperator = "not_contains"

	// OperatorNotEquals 不等于：字段值不等于指定值
	OperatorNotEquals ConditionOperator = "not_equals"

	// OperatorLess 小于：字段值小于指定值
	// 如：金额 less "100"
	OperatorLess ConditionOperator = "less"

	// OperatorMore 大于：字段值大于指定值
	// 如：金额 more "1000"
	OperatorMore ConditionOperator = "more"

	// OperatorIsEmpty 为空：字段值为空
	OperatorIsEmpty ConditionOperator = "is_empty"

	// OperatorIsNotEmpty 不为空：字段值不为空
	OperatorIsNotEmpty ConditionOperator = "is_not_empty"
)

// RuleCondition 规则条件模型，对应 rule_conditions 表
//
// 功能说明：
// - 定义规则的匹配条件
// - 每条规则可以有多个条件，条件之间通过 LogicType（AND/OR）组合
// - 条件由三部分组成：字段（Field）+ 运算符（Operator）+ 值（Value）
//
// 条件匹配示例：
//   描述 contains "星巴克" -> 匹配描述中包含"星巴克"的交易
//   金额 more "1000" -> 匹配金额大于 1000 的交易
//   源账户 equals "3" -> 匹配源账户 ID 为 3 的交易
type RuleCondition struct {
	// ID 条件唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// RuleID 所属规则 ID
	// 关联到 rules 表
	// gorm:"index" 创建索引，加速按规则查询条件
	RuleID uint64 `gorm:"not null;index" json:"rule_id"`

	// Field 条件字段
	// 取值为 ConditionField 枚举
	// 指定要匹配的交易字段
	// gorm:"size:20" 限制最大长度为 20 个字符
	Field ConditionField `gorm:"not null;size:20" json:"field"`

	// Operator 条件运算符
	// 取值为 ConditionOperator 枚举
	// 指定字段与值之间的比较方式
	// gorm:"size:20" 限制最大长度为 20 个字符
	Operator ConditionOperator `gorm:"not null;size:20" json:"operator"`

	// Value 条件值
	// 与 Field 和 Operator 配合使用
	// 如：Field=description, Operator=contains, Value="星巴克"
	// gorm:"size:255" 限制最大长度为 255 个字符
	Value string `gorm:"not null;size:255" json:"value"`

	// CreatedAt 条件创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 RuleCondition 模型对应的数据库表名为 rule_conditions
func (RuleCondition) TableName() string { return "rule_conditions" }

// ActionType 操作类型枚举
// 定义规则条件满足时可执行的操作类型
type ActionType string

const (
	// ActionTypeSetCategory 设置分类
	// 将交易的分类设置为指定值
	// Value: 分类 ID
	ActionTypeSetCategory ActionType = "set_category"

	// ActionTypeAddTag 添加标签
	// 为交易添加指定标签
	// Value: 标签名称
	ActionTypeAddTag ActionType = "add_tag"

	// ActionTypeSetNotes 设置备注
	// 将交易的备注设置为指定值（覆盖原有备注）
	// Value: 备注内容
	ActionTypeSetNotes ActionType = "set_notes"

	// ActionTypeSetBudget 设置预算
	// 将交易关联到指定预算
	// Value: 预算 ID
	ActionTypeSetBudget ActionType = "set_budget"

	// ActionTypeRemoveTag 移除标签
	// 从交易中移除指定标签
	// Value: 标签名称
	ActionTypeRemoveTag ActionType = "remove_tag"

	// ActionTypeSetDescription 设置描述
	// 将交易的描述设置为指定值（覆盖原有描述）
	// Value: 描述内容
	ActionTypeSetDescription ActionType = "set_description"

	// ActionTypeClearCategory 清除分类
	// 移除交易的分类关联
	// Value: 无需值
	ActionTypeClearCategory ActionType = "clear_category"

	// ActionTypeClearBudget 清除预算
	// 移除交易的预算关联
	// Value: 无需值
	ActionTypeClearBudget ActionType = "clear_budget"

	// ActionTypeClearNotes 清除备注
	// 清空交易的备注
	// Value: 无需值
	ActionTypeClearNotes ActionType = "clear_notes"

	// ActionTypeAppendNotes 追加备注
	// 在交易原有备注后追加内容
	// Value: 追加的备注内容
	ActionTypeAppendNotes ActionType = "append_notes"

	// ActionTypePrependNotes 前置备注
	// 在交易原有备注前插入内容
	// Value: 前置的备注内容
	ActionTypePrependNotes ActionType = "prepend_notes"
)

// RuleAction 规则操作模型，对应 rule_actions 表
//
// 功能说明：
// - 定义规则条件满足时执行的操作
// - 每条规则可以有多个操作，条件满足时依次执行
// - 操作由两部分组成：操作类型（Type）+ 操作值（Value）
//
// 操作执行示例：
//   Type=set_category, Value="5" -> 将交易分类设置为 ID=5 的分类
//   Type=add_tag, Value="咖啡" -> 为交易添加"咖啡"标签
//   Type=append_notes, Value="自动标记" -> 在备注后追加"自动标记"
type RuleAction struct {
	// ID 操作唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// RuleID 所属规则 ID
	// 关联到 rules 表
	// gorm:"index" 创建索引，加速按规则查询操作
	RuleID uint64 `gorm:"not null;index" json:"rule_id"`

	// Type 操作类型
	// 取值为 ActionType 枚举
	// 指定要执行的操作类型
	// gorm:"size:20" 限制最大长度为 20 个字符
	Type ActionType `gorm:"not null;size:20" json:"type"`

	// Value 操作值
	// 与 Type 配合使用，指定操作的具体参数
	// 如：Type=set_category, Value="5"（分类 ID）
	// 如：Type=add_tag, Value="咖啡"（标签名称）
	// gorm:"size:255" 限制最大长度为 255 个字符
	Value string `gorm:"not null;size:255" json:"value"`

	// CreatedAt 操作创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 RuleAction 模型对应的数据库表名为 rule_actions
func (RuleAction) TableName() string { return "rule_actions" }
