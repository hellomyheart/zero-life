// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型，对应categories表
// 支持两级分类结构，用于对交易进行分类管理
type Category struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`           // 所属用户ID
	Name      string         `gorm:"not null;size:100" json:"name"`           // 分类名称
	ParentID  *uint64        `gorm:"index" json:"parent_id"`                  // 父分类ID，nil表示顶级分类
	Icon      string         `gorm:"size:50" json:"icon"`                     // 图标标识
	Notes     string         `gorm:"type:text" json:"notes"`                  // 备注信息
	SortOrder int           `gorm:"default:0" json:"sort_order"`             // 排序序号
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间

	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子分类列表
}

// TableName 指定表名
func (Category) TableName() string { return "categories" }
