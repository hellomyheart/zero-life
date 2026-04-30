package repository

import (
	"errors"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBudgetRepository(readDB, writeDB *gorm.DB) *BudgetRepository {
	return &BudgetRepository{readDB: readDB, writeDB: writeDB}
}

func (r *BudgetRepository) WriteDB() *gorm.DB {
	return r.writeDB
}

func (r *BudgetRepository) Create(budget *model.Budget, categoryIDs []uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
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
	return r.GetByIDWithDB(r.readDB, id, userID)
}

func (r *BudgetRepository) GetByIDWithDB(db *gorm.DB, id, userID uint64) (*model.Budget, error) {
	var budget model.Budget
	if err := db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Categories").
		First(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) List(userID uint64) ([]model.Budget, error) {
	var budgets []model.Budget
	if err := r.readDB.Where("user_id = ?", userID).
		Preload("Categories").
		Order("id ASC").Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) Update(budget *model.Budget, categoryIDs []uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(budget).Error; err != nil {
			return err
		}
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
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("budget_id = ?", id).Delete(&model.BudgetCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("budget_id = ?", id).Delete(&model.BudgetHistory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Budget{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *BudgetRepository) GetHistory(budgetID uint64) ([]model.BudgetHistory, error) {
	var history []model.BudgetHistory
	if err := r.readDB.Where("budget_id = ?", budgetID).Order("period_start DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *BudgetRepository) CreateHistory(history *model.BudgetHistory) error {
	return r.writeDB.Create(history).Error
}

func (r *BudgetRepository) ListAllEnabled() ([]model.Budget, error) {
	var budgets []model.Budget
	if err := r.readDB.Where("is_enabled = ?", true).
		Preload("Categories").
		Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) UpsertHistory(history *model.BudgetHistory) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		var existing model.BudgetHistory
		err := tx.Where("budget_id = ? AND period_start = ?", history.BudgetID, history.PeriodStart).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(history).Error
			}
			return err
		}
		return tx.Model(&existing).Updates(map[string]interface{}{
			"period_end": history.PeriodEnd,
			"spent":      history.Spent,
			"updated_at": history.UpdatedAt,
		}).Error
	})
}
