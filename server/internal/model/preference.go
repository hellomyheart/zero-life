package model

import "time"

type Preference struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;uniqueIndex:idx_user_name" json:"user_id"`
	Name      string    `gorm:"not null;size:100;uniqueIndex:idx_user_name" json:"name"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (Preference) TableName() string { return "preferences" }
