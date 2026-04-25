package repository

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(budget *model.Budget, categoryIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(budget).Error; err != nil {
			return err
		}
		if len(categoryIDs) > 0 {
			budgetCategories := make([]model.BudgetCategory, 0, len(categoryIDs))
			for _, catID := range categoryIDs {
				budgetCategories = append(budgetCategories, model.BudgetCategory{
					BudgetID:   budget.ID,
					CategoryID: catID,
				})
			}
			if err := tx.Create(&budgetCategories).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BudgetRepository) GetByID(id, userID uint64) (*model.Budget, error) {
	var budget model.Budget
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Categories").
		First(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) List(userID uint64) ([]model.Budget, error) {
	var budgets []model.Budget
	if err := r.db.Where("user_id = ?", userID).
		Preload("Categories").
		Order("id ASC").Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) Update(budget *model.Budget, categoryIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(budget).Error; err != nil {
			return err
		}
		// Replace categories: delete old, insert new
		if err := tx.Where("budget_id = ?", budget.ID).Delete(&model.BudgetCategory{}).Error; err != nil {
			return err
		}
		if len(categoryIDs) > 0 {
			budgetCategories := make([]model.BudgetCategory, 0, len(categoryIDs))
			for _, catID := range categoryIDs {
				budgetCategories = append(budgetCategories, model.BudgetCategory{
					BudgetID:   budget.ID,
					CategoryID: catID,
				})
			}
			if err := tx.Create(&budgetCategories).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BudgetRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete budget categories
		if err := tx.Where("budget_id = ?", id).Delete(&model.BudgetCategory{}).Error; err != nil {
			return err
		}
		// Delete budget history
		if err := tx.Where("budget_id = ?", id).Delete(&model.BudgetHistory{}).Error; err != nil {
			return err
		}
		// Delete the budget
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Budget{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *BudgetRepository) GetHistory(budgetID uint64) ([]model.BudgetHistory, error) {
	var history []model.BudgetHistory
	if err := r.db.Where("budget_id = ?", budgetID).Order("period_start DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *BudgetRepository) CreateHistory(history *model.BudgetHistory) error {
	return r.db.Create(history).Error
}

// GetSpentAmount 获取预算已花费金额
func (r *BudgetRepository) GetSpentAmount(budgetID, userID uint64) (decimal.Decimal, error) {
	// 获取预算信息（包含关联的分类）
	budget, err := r.GetByID(budgetID, userID)
	if err != nil {
		return decimal.Zero, err
	}

	// 计算当前周期的花费
	query := r.db.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "withdrawal")

	// 如果预算关联了分类，只统计这些分类的交易
	// Categories 是 []Category 类型，cat.ID 就是分类的ID
	if len(budget.Categories) > 0 {
		categoryIDs := make([]uint64, len(budget.Categories))
		for i, cat := range budget.Categories {
			categoryIDs[i] = cat.ID
		}
		query = query.Where("category_id IN ?", categoryIDs)
	}

	// 根据预算周期计算时间范围
	now := time.Now()
	var startDate, endDate time.Time
	switch budget.Period {
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	case "yearly":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0).Add(-time.Second)
	default:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	}

	query = query.Where("date >= ? AND date <= ?", startDate, endDate)

	// 查询总金额
	var result struct {
		Total decimal.Decimal
	}
	if err := query.Select("COALESCE(SUM(amount), 0) as total").Scan(&result).Error; err != nil {
		return decimal.Zero, err
	}

	return result.Total, nil
}
