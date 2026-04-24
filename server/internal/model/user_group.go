// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserGroup 用户组模型，对应user_groups表
// 用于组织和管理用户，支持多租户场景
type UserGroup struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`         // 用户组ID，主键自增
	Title     string         `gorm:"not null;size:255" json:"title"`             // 用户组标题
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`                 // 创建时间
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`                 // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                             // 软删除时间
}

// TableName 指定表名
func (UserGroup) TableName() string { return "user_groups" }

// UserGroupMember 用户组成员关联表模型，对应user_group_members表
// 记录用户与用户组的关联关系
type UserGroupMember struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`          // 关联ID，主键自增
	UserGroupID uint64         `gorm:"not null;index" json:"user_group_id"`         // 用户组ID
	UserID      uint64         `gorm:"not null;index" json:"user_id"`               // 用户ID
	IsOwner     bool           `gorm:"default:false" json:"is_owner"`               // 是否为组所有者
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`                  // 创建时间
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`                  // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                              // 软删除时间

	UserGroup UserGroup `gorm:"foreignKey:UserGroupID" json:"user_group,omitempty"` // 用户组信息
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`             // 用户信息
}

// TableName 指定表名
func (UserGroupMember) TableName() string { return "user_group_members" }
