package model

import (
	"time"

	"gorm.io/gorm"
)

type TransactionLink struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	LinkTypeID    uint64         `gorm:"not null;index" json:"link_type_id"`
	SourceID      uint64         `gorm:"not null;index" json:"source_id"`
	DestinationID uint64         `gorm:"not null;index" json:"destination_id"`
	Comment       string         `gorm:"size:500" json:"comment"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	LinkType      LinkType       `gorm:"foreignKey:LinkTypeID" json:"link_type,omitempty"`
}

func (TransactionLink) TableName() string { return "transaction_links" }
