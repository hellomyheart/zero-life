package model

import (
	"time"

	"gorm.io/gorm"
)

type LinkType struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Outward       string         `gorm:"not null;size:100" json:"outward"`
	Inward        string         `gorm:"not null;size:100" json:"inward"`
	IsDirectional bool           `gorm:"default:false" json:"is_directional"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LinkType) TableName() string { return "link_types" }
