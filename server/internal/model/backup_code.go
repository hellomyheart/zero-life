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
