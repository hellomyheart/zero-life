package model

import "time"

type Configuration struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (Configuration) TableName() string { return "configurations" }
