// Package model 数据模型层，对应数据库结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// ObjectGroup 对象分组模型，对应object_groups表
// 管理各类实体（账户、分类、标签等）的分组和排序
// 采用多态设计：通过GroupableType+GroupableID关联不同类型的实体
// 例如：将多个账户归为"银行账户"组，将多个分类归为"日常支出"组
type ObjectGroup struct {
	// ID 分组唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户ID，每个用户有自己独立的分组
	// gorm:"index" 创建索引，加速按用户查询分组
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Name 分组名称，如"银行账户"、"投资账户"等
	// gorm:"size:100" 限制名称最大长度为100个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// GroupableType 关联对象的类型（多态），如"account"、"category"等
	// 与GroupableID配合实现多态关联，类似于面向对象中的多态
	// gorm:"size:50" 限制类型名称最大长度为50个字符
	GroupableType string `gorm:"not null;size:50" json:"groupable_type"`

	// GroupableID 关联对象的ID（多态），与GroupableType配合使用
	// 例如：GroupableType="account" + GroupableID=3 表示关联ID为3的账户
	// gorm:"index" 创建索引，加速按对象ID查询分组
	GroupableID uint64 `gorm:"not null;index" json:"groupable_id"`

	// CreatedAt 分组创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 分组最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间，GORM软删除字段
	// gorm.DeletedAt 是GORM内置的软删除类型，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM查询时会自动过滤已软删除的记录
	// json:"-" 表示JSON序列化时忽略此字段，不返回给前端
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名为object_groups
func (ObjectGroup) TableName() string { return "object_groups" }
