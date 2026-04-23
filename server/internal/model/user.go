package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	Nickname  string         `gorm:"not null;size:100" json:"nickname"`
	Language  string         `gorm:"size:10;default:zh-CN" json:"language"`
	Timezone  string         `gorm:"size:50;default:Asia/Shanghai" json:"timezone"`
	Role      string         `gorm:"size:20;default:user" json:"role"`
	MFASecret string         `gorm:"size:255" json:"-"`
	MFAEnabled bool          `gorm:"default:false" json:"mfa_enabled"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
