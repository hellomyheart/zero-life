// Package model 数据模型层，对应数据库结构
// ObjectGroup 对象分组模型，管理实体分组排序
package model

import (
	"time"

	"gorm.io/gorm"
)

type ObjectGroup struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64         `gorm:"not null;index" json:"user_id"`
	Name          string         `gorm:"not null;size:100" json:"name"`
	GroupableType string         `gorm:"not null;size:50" json:"groupable_type"`
	GroupableID   uint64         `gorm:"not null;index" json:"groupable_id"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ObjectGroup) TableName() string { return "object_groups" }