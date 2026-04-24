// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"gorm.io/gorm"
)

// LinkType 链接类型模型，对应 link_types 表
// 
// 功能说明：
// - 定义交易之间的关联类型
// - 支持双向描述（正向/反向）
// - 支持有向/无向关联
//
// 使用场景：
// - 退款关联：原始交易 <-> 退款交易
// - 分期关联：首付交易 <-> 分期交易
// - 冲销关联：错误交易 <-> 冲销交易
// - 对账关联：系统交易 <-> 银行交易
//
// 双向描述：
// - Outward: 从 A 到 B 的描述（如"退款"）
// - Inward: 从 B 到 A 的描述（如"被退款"）
// - IsDirectional: 是否区分方向
//
// 示例 1（有向）：
//   LinkType: "退款"
//   Outward: "退款"（交易 A 退款到交易 B）
//   Inward: "被退款"（交易 B 被交易 A 退款）
//   IsDirectional: true
//
// 示例 2（无向）：
//   LinkType: "关联"
//   Outward: "关联到"
//   Inward: "关联到"
//   IsDirectional: false
//
// 预定义链接类型：
// - related: 关联（通用）
// - reconciled: 已对账
// - rolled_back: 已回滚/冲销
// - refund: 退款
// - partial_refund: 部分退款
// - installment: 分期
//
// 与 TransactionLink 的关系：
// - LinkType 定义关联的"类型"
// - TransactionLink 使用 LinkType 创建具体的"关联"
type LinkType struct {
	// ID 链接类型唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// Name 链接类型名称
	// 如："退款"、"关联"、"冲销"
	// gorm:"uniqueIndex" 创建唯一索引，确保名称不重复
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"uniqueIndex;not null;size:100" json:"name"`
	
	// Outward 正向描述
	// 从源交易到目标交易的描述
	// 如："退款"、"关联到"、"冲销"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Outward string `gorm:"not null;size:100" json:"outward"`
	
	// Inward 反向描述
	// 从目标交易到源交易的描述
	// 如："被退款"、"关联自"、"被冲销"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Inward string `gorm:"not null;size:100" json:"inward"`
	
	// IsDirectional 是否有方向性
	// true: 区分方向，Outward != Inward
	// false: 不区分方向，Outward == Inward
	// 
	// 示例：
	//   "退款" 是有向的：A 退款给 B，但 B 没有退款给 A
	//   "关联" 是无向的：A 关联到 B，B 也关联到 A
	// 
	// 默认值：false
	IsDirectional bool `gorm:"default:false" json:"is_directional"`
	
	// CreatedAt 链接类型创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 链接类型最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 LinkType 模型对应的数据库表名为 link_types
func (LinkType) TableName() string { return "link_types" }
