// Package model 数据模型层，对应数据库结构
// BackupCode MFA备用码模型，存储两步验证备用码
package model

import "time"

type BackupCode struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;index" json:"user_id"`
	Code      string     `gorm:"not null;size:20;index" json:"code"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
}

func (BackupCode) TableName() string { return "backup_codes" }
