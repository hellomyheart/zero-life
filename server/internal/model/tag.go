// Package model 定义了系统的数据模型，对应数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型，对应tags表
// 用于对交易进行标记和分组，支持自定义颜色
type Tag struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`       // 所属用户ID
	Name      string         `gorm:"not null;size:100" json:"name"`       // 标签名称
	Color     string         `gorm:"size:7;default:#409EFF" json:"color"` // 标签颜色，十六进制格式
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                      // 软删除时间
}

// TableName 指定表名
func (Tag) TableName() string { return "tags" }

// TransactionTag 交易-标签关联模型，对应transaction_tags表
// 多对多关联的中间表
type TransactionTag struct {
	TransactionID uint64 `gorm:"primaryKey" json:"transaction_id"` // 交易ID
	TagID         uint64 `gorm:"primaryKey" json:"tag_id"`         // 标签ID
}

// TableName 指定表名
func (TransactionTag) TableName() string { return "transaction_tags" }
