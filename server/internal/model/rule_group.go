package model

import (
	"time"

	"gorm.io/gorm"
)

type RuleGroup struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Order     int            `gorm:"default:0" json:"order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RuleGroup) TableName() string { return "rule_groups" }
