package model

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64         `gorm:"not null;index" json:"user_id"`
	AttachableType string         `gorm:"not null;size:50;index" json:"attachable_type"`
	AttachableID   uint64         `gorm:"not null;index" json:"attachable_id"`
	Filename       string         `gorm:"not null;size:255" json:"filename"`
	Mime           string         `gorm:"not null;size:255" json:"mime"`
	Size           int64          `gorm:"not null" json:"size"`
	Path           string         `gorm:"not null;size:500" json:"path"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Attachment) TableName() string { return "attachments" }
