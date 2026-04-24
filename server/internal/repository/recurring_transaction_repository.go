package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type RecurringTransactionRepository struct {
	db *gorm.DB
}

func NewRecurringTransactionRepository(db *gorm.DB) *RecurringTransactionRepository {
	return &RecurringTransactionRepository{db: db}
}

func (r *RecurringTransactionRepository) Create(rt *model.RecurringTransaction) error {
	return r.db.Create(rt).Error
}

func (r *RecurringTransactionRepository) GetByID(id, userID uint64) (*model.RecurringTransaction, error) {
	var rt model.RecurringTransaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RecurringTransactionRepository) List(userID uint64, active *bool, offset, limit int) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	query := r.db.Where("user_id = ?", userID)
	if active != nil {
		query = query.Where("is_active = ?", *active)
	}
	if err := query.Order("next_occurrence ASC").Offset(offset).Limit(limit).Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

func (r *RecurringTransactionRepository) Count(userID uint64, active *bool) (int64, error) {
	var count int64
	query := r.db.Model(&model.RecurringTransaction{}).Where("user_id = ?", userID)
	if active != nil {
		query = query.Where("is_active = ?", *active)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RecurringTransactionRepository) Update(rt *model.RecurringTransaction) error {
	return r.db.Save(rt).Error
}

func (r *RecurringTransactionRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("recurring_transaction_id = ?", id).Delete(&model.RecurringTransactionLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecurringTransaction{}).Error
	})
}

func (r *RecurringTransactionRepository) GetDueRecurringTransactions(userID uint64) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := gorm.Expr("NOW()")
	if err := r.db.Where("user_id = ? AND is_active = ? AND next_occurrence <= ?", userID, true, now).
		Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

func (r *RecurringTransactionRepository) CreateLog(log *model.RecurringTransactionLog) error {
	return r.db.Create(log).Error
}

func (r *RecurringTransactionRepository) ListLogs(recurringTransactionID uint64) ([]model.RecurringTransactionLog, error) {
	var logs []model.RecurringTransactionLog
	if err := r.db.Where("recurring_transaction_id = ?", recurringTransactionID).
		Order("occurrence_date DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
