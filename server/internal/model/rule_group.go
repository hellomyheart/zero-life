// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单、规则等
package model

import (
	"time"

	"gorm.io/gorm"
)

// RuleGroup 规则组模型，对应 rule_groups 表
// 
// 功能说明：
// - 将规则分组管理，便于组织和执行
// - 按组顺序依次执行规则（Order 字段控制）
// - 可以整体启用/禁用一组规则
//
// 使用场景：
// - 按业务场景分组：如"信用卡规则组"、"报销规则组"
// - 按执行时机分组：如"导入时规则"、"手动执行规则"
// - 按优先级分组：如"高优先级组"、"常规组"
//
// 执行顺序：
// 1. 系统按 Order 字段升序排列规则组
// 2. 依次执行每个组内的规则
// 3. 组内规则按 Priority 字段排序
// 4. 如果组被禁用（IsActive=false），跳过该组所有规则
//
// 示例：
//   规则组 1: Order=1, 名称="导入时自动处理"
//     - 规则 1: 自动设置分类
//     - 规则 2: 自动打标签
//   规则组 2: Order=2, 名称="月末批量处理"
//     - 规则 3: 关联账单
//     - 规则 4: 更新预算
type RuleGroup struct {
	// ID 规则组唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID
	// 每个用户有自己的规则组集合，用户间数据隔离
	// gorm:"index" 加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// Name 规则组名称
	// 如："导入时自动处理"、"信用卡规则"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`
	
	// Order 排序序号
	// 数值越小越先执行
	// 默认值：0
	// 
	// 执行顺序示例：
	//   Order=1 的组 -> Order=10 的组 -> Order=100 的组
	Order int `gorm:"default:0" json:"order"`
	
	// IsActive 是否激活
	// true: 正常执行该组规则
	// false: 跳过该组所有规则（但保留配置）
	// 默认值：true
	IsActive bool `gorm:"default:true" json:"is_active"`
	
	// CreatedAt 规则组创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 规则组最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 RuleGroup 模型对应的数据库表名为 rule_groups
func (RuleGroup) TableName() string { return "rule_groups" }
