// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// RuleGroup 规则组模型，对应rule_groups表
// 将规则分组管理，按组顺序依次执行
type RuleGroup struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`  // 所属用户ID
	Name      string         `gorm:"not null;size:100" json:"name"`  // 规则组名称
	Order     int            `gorm:"default:0" json:"order"`         // 排序序号，数值越小越先执行
	IsActive  bool           `gorm:"default:true" json:"is_active"`  // 是否激活
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                 // 软删除时间
}

// TableName 指定表名
func (RuleGroup) TableName() string { return "rule_groups" }
