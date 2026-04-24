// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// LinkType 链接类型模型，对应link_types表
// 定义交易之间的关联类型，支持双向描述（如"退款"/"被退款"）
type LinkType struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"uniqueIndex;not null;size:100" json:"name"`  // 链接类型名称
	Outward       string         `gorm:"not null;size:100" json:"outward"`           // 正向描述（如"退款"）
	Inward        string         `gorm:"not null;size:100" json:"inward"`            // 反向描述（如"被退款"）
	IsDirectional bool           `gorm:"default:false" json:"is_directional"`        // 是否有方向性
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                             // 软删除时间
}

// TableName 指定表名
func (LinkType) TableName() string { return "link_types" }
