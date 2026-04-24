// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// Attachment 附件模型，对应attachments表
// 支持为交易、账户等对象附加文件（图片、PDF等）
type Attachment struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64         `gorm:"not null;index" json:"user_id"`                // 所属用户ID
	AttachableType string         `gorm:"not null;size:50;index" json:"attachable_type"` // 关联对象类型（transaction/account等）
	AttachableID   uint64         `gorm:"not null;index" json:"attachable_id"`           // 关联对象ID
	Filename       string         `gorm:"not null;size:255" json:"filename"`             // 原始文件名
	Mime           string         `gorm:"not null;size:255" json:"mime"`                 // MIME类型
	Size           int64          `gorm:"not null" json:"size"`                          // 文件大小（字节）
	Path           string         `gorm:"not null;size:500" json:"path"`                 // 存储路径
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                                // 软删除时间
}

// TableName 指定表名
func (Attachment) TableName() string { return "attachments" }
