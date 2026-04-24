package model

import (
	"time"

	"gorm.io/gorm"
)

type Preference struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;uniqueIndex:idx_pref_user_key" json:"user_id"`
	Key       string         `gorm:"not null;size:100;uniqueIndex:idx_pref_user_key" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Preference) TableName() string { return "preferences" }