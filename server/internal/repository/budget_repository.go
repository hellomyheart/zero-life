package repository

import (
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
