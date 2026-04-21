package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetPeriod string

const (
	BudgetPeriodMonthly BudgetPeriod = "monthly"
	BudgetPeriodYearly  BudgetPeriod = "yearly"
)

type Budget struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Amount    decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	Period    BudgetPeriod   `gorm:"type:enum('monthly','yearly');not null" json:"period"`
	IsEnabled bool           `gorm:"default:true" json:"is_enabled"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Categories []Category `gorm:"many2many:budget_categories" json:"categories,omitempty"`
}

func (Budget) TableName() string { return "budgets" }

type BudgetCategory struct {
	BudgetID   uint64 `gorm:"primaryKey" json:"budget_id"`
	CategoryID uint64 `gorm:"primaryKey" json:"category_id"`
}

func (BudgetCategory) TableName() string { return "budget_categories" }

type BudgetHistory struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BudgetID    uint64          `gorm:"not null;index" json:"budget_id"`
	PeriodStart time.Time       `gorm:"not null" json:"period_start"`
	PeriodEnd   time.Time       `gorm:"not null" json:"period_end"`
	Amount      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	Spent       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"spent"`
	CreatedAt   time.Time       `gorm:"not null" json:"created_at"`
}

func (BudgetHistory) TableName() string { return "budget_history" }
