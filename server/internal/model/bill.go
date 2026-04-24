// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RepeatRule 账单重复规则枚举
// 
// 定义账单的重复周期，用于计算下次到期日
// 
// 使用场景：
// - Daily: 每日重复（如：每日餐费）
// - Weekly: 每周重复（如：每周保洁）
// - Monthly: 每月重复（如：房租、网费）
// - Yearly: 每年重复（如：保险费、订阅费）
//
// 示例：
//   房租账单：RepeatRule = "monthly", NextDue = 2026-05-01
//   到期后自动计算下次到期日：2026-06-01
type RepeatRule string

const (
	// RepeatRuleDaily 每日重复
	RepeatRuleDaily   RepeatRule = "daily"
	
	// RepeatRuleWeekly 每周重复
	RepeatRuleWeekly  RepeatRule = "weekly"
	
	// RepeatRuleMonthly 每月重复
	// 最常见的账单类型：房租、网费、手机费、订阅服务
	RepeatRuleMonthly RepeatRule = "monthly"
	
	// RepeatRuleYearly 每年重复
	// 年度费用：保险费、会员费、年检费
	RepeatRuleYearly  RepeatRule = "yearly"
)

// Bill 账单模型，对应 bills 表
// 
// 功能说明：
// - 管理周期性账单，自动提醒到期付款
// - 支持多种重复规则（每日/每周/每月/每年）
// - 可关联支出账户和分类
// - 自动计算下次到期日
//
// 使用场景：
// - 固定支出管理：房租、房贷、车贷
// - 订阅服务跟踪：Netflix、Spotify、软件订阅
// - 定期账单提醒：水电煤、物业费、网费
// - 保险费用管理：车险、健康险、寿险
//
// 自动匹配逻辑：
// 1. 系统定期扫描到期账单（NextDue <= 当前日期）
// 2. 根据账单名称、金额、账户自动匹配交易
// 3. 匹配成功后标记账单为"已支付"
// 4. 计算下次到期日（根据 RepeatRule）
//
// 与交易的关系：
// - 账单是"计划"，交易是"实际"
// - 一笔账单可以对应多笔交易（分期支付）
// - 交易可以通过 BillID 关联到账单
//
// 示例：
//   账单：房租 - 5000 元/月 - 每月 1 日到期
//   交易：5 月 1 日支付房租 5000 元 -> 关联到账单
type Bill struct {
	// ID 账单唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID
	// 每个用户有自己的账单集合，用户间数据隔离
	// gorm:"index" 加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// Name 账单名称
	// 如："房租"、"网费"、"Netflix 订阅"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`
	
	// Amount 账单金额
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 +4 位小数
	Amount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	
	// RepeatRule 重复规则
	// 决定账单的重复周期：daily/weekly/monthly/yearly
	// gorm:"size:20" 限制最大长度为 20 个字符
	RepeatRule RepeatRule `gorm:"not null;size:20" json:"repeat_rule"`
	
	// NextDue 下次到期日期
	// 账单的下一个付款截止日
	// 到期后系统会自动提醒，并在支付后计算下一次到期日
	NextDue time.Time `gorm:"not null" json:"next_due"`
	
	// SourceID 支出账户 ID（可选）
	// 关联的资产账户，用于支付该账单
	// 如：用"招商银行卡"支付房租
	// json:"source_id" 允许为空
	SourceID *uint64 `json:"source_id"`
	
	// CategoryID 分类 ID（可选）
	// 关联的支出分类
	// 如：房租账单关联到"居住"分类
	// json:"category_id" 允许为空
	CategoryID *uint64 `json:"category_id"`
	
	// Notes 备注信息
	// 记录账单的附加说明，如房东联系方式、缴费账号等
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`
	
	// CreatedAt 账单创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 账单最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 Bill 模型对应的数据库表名为 bills
func (Bill) TableName() string { return "bills" }
