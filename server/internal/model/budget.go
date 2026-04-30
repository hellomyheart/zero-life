// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BudgetPeriod 预算周期类型
// 定义预算的时间范围，用于控制支出上限的统计周期
type BudgetPeriod string

const (
	BudgetPeriodDaily     BudgetPeriod = "daily"
	BudgetPeriodWeekly    BudgetPeriod = "weekly"
	BudgetPeriodMonthly   BudgetPeriod = "monthly"
	BudgetPeriodQuarterly BudgetPeriod = "quarterly"
	BudgetPeriodYearly    BudgetPeriod = "yearly"
)

// Budget 预算模型，对应 budgets 表
//
// 功能说明：
// - 设定分类支出上限，跟踪预算使用情况
// - 支持日度、周度、月度、季度、年度五种预算周期
// - 一个预算可以关联多个分类（多对多关系）
// - 通过预算历史记录跟踪每个周期的执行情况
//
// 使用场景：
// - 月度餐饮预算：设定每月餐饮支出上限 3000 元
// - 年度旅游预算：设定全年旅游支出上限 20000 元
// - 综合预算：将"餐饮+交通+购物"合并为一个预算
//
// 与其他模型的关系：
// - Category: 多对多关系，一个预算可覆盖多个分类
// - BudgetHistory: 一对多关系，记录每个预算周期的执行情况
// - Transaction: 间接关联，通过分类统计交易金额
//
// 预算计算逻辑：
//   已用金额 = 当前周期内，关联分类下所有支出交易的金额总和
//   剩余金额 = Amount - 已用金额
//   使用率 = 已用金额 / Amount × 100%
//
// 示例：
//   预算：月度餐饮预算
//   金额：3000 元
//   关联分类：[餐饮, 零食, 咖啡]
//   本月已用：2500 元
//   剩余：500 元
//   使用率：83.3%
type Budget struct {
	// ID 预算唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的预算集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询预算
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Name 预算名称
	// 用户自定义的预算名称，便于识别
	// 如："月度餐饮预算"、"年度旅游预算"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// Amount 预算金额上限
	// 当前周期内允许的最大支出金额
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// Period 预算周期
	// 取值为 BudgetPeriod 枚举：daily/weekly/monthly/quarterly/yearly
	// 决定预算的重置频率
	// gorm:"size:20" 限制最大长度为 20 个字符
	Period BudgetPeriod `gorm:"not null;size:20" json:"period"`

	// IsEnabled 是否启用
	// true: 预算生效，系统会跟踪支出并提醒
	// false: 预算暂停，不跟踪支出
	// 默认值：true
	IsEnabled bool `gorm:"default:true" json:"is_enabled"`

	// CreatedAt 预算创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 预算最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Categories 关联的分类列表（多对多关系）
	// 通过 budget_categories 中间表实现
	// 一个预算可以覆盖多个分类的支出
	// json:"categories,omitempty" 当字段为零值时 JSON 序列化时忽略
	Categories []Category `gorm:"many2many:budget_categories" json:"categories,omitempty"`
}

// TableName 指定 Budget 模型对应的数据库表名为 budgets
func (Budget) TableName() string { return "budgets" }

// BudgetCategory 预算-分类关联模型，对应 budget_categories 表
//
// 功能说明：
// - 实现预算与分类的多对多关联
// - 一个预算可以关联多个分类
// - 一个分类可以被多个预算覆盖
//
// 数据结构：
// - 复合主键：(BudgetID, CategoryID)
// - 无额外字段，纯关联表
//
// 示例：
//   预算"月度餐饮预算"(ID=1) 关联分类：餐饮(ID=5)、零食(ID=6)、咖啡(ID=7)
//   即：{1,5}, {1,6}, {1,7}
type BudgetCategory struct {
	// BudgetID 预算 ID，复合主键的一部分
	// 关联到 budgets 表的 ID 字段
	BudgetID uint64 `gorm:"primaryKey" json:"budget_id"`

	// CategoryID 分类 ID，复合主键的另一部分
	// 关联到 categories 表的 ID 字段
	CategoryID uint64 `gorm:"primaryKey" json:"category_id"`
}

// TableName 指定 BudgetCategory 模型对应的数据库表名为 budget_categories
func (BudgetCategory) TableName() string { return "budget_categories" }

// BudgetHistory 预算历史模型，对应 budget_history 表
//
// 功能说明：
// - 记录每个预算周期的实际支出和预算金额
// - 用于生成预算执行报告和趋势分析
// - 每个预算周期结束时自动生成一条记录
//
// 使用场景：
// - 查看历史预算执行情况
// - 分析支出趋势（逐月对比）
// - 评估预算设定是否合理
//
// 示例：
//   预算：月度餐饮预算 3000 元
//   2026年4月：预算 3000 元，实际支出 2800 元（节省 200 元）
//   2026年3月：预算 3000 元，实际支出 3200 元（超支 200 元）
type BudgetHistory struct {
	// ID 预算历史记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// BudgetID 预算 ID
	// 关联到 budgets 表
	BudgetID uint64 `gorm:"not null;uniqueIndex:idx_budget_period" json:"budget_id"`

	// PeriodStart 周期开始日期
	// 如月度预算：2026-04-01
	// 如年度预算：2026-01-01
	PeriodStart time.Time `gorm:"not null;uniqueIndex:idx_budget_period" json:"period_start"`

	// PeriodEnd 周期结束日期
	// 如月度预算：2026-04-30
	// 如年度预算：2026-12-31
	PeriodEnd time.Time `gorm:"not null" json:"period_end"`

	// Amount 预算金额
	// 该周期设定的预算上限
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`

	// Spent 实际支出金额
	// 该周期内关联分类的支出总和
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 + 4 位小数
	Spent decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"spent"`

	// CreatedAt 历史记录创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 历史记录最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定 BudgetHistory 模型对应的数据库表名为 budget_history
func (BudgetHistory) TableName() string { return "budget_history" }
