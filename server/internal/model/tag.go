package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Color     string         `gorm:"size:7;default:#409EFF" json:"color"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tag) TableName() string { return "tags" }

type TransactionTag struct {
	TransactionID uint64 `gorm:"primaryKey" json:"transaction_id"`
	TagID         uint64 `gorm:"primaryKey" json:"tag_id"`
}

func (TransactionTag) TableName() string { return "transaction_tags" }
