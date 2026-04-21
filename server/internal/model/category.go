package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	ParentID  *uint64        `gorm:"index" json:"parent_id"`
	Icon      string         `gorm:"size:50" json:"icon"`
	Notes     string         `gorm:"type:text" json:"notes"`
	SortOrder int           `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Category) TableName() string { return "categories" }
