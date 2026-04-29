package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type RecurringTransactionRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewRecurringTransactionRepository(readDB, writeDB *gorm.DB) *RecurringTransactionRepository {
	return &RecurringTransactionRepository{readDB: readDB, writeDB: writeDB}
}

func (r *RecurringTransactionRepository) Create(rt *model.RecurringTransaction) error {
	return r.writeDB.Create(rt).Error
}

func (r *RecurringTransactionRepository) GetByID(id, userID uint64) (*model.RecurringTransaction, error) {
	var rt model.RecurringTransaction
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RecurringTransactionRepository) List(userID uint64, active *bool, offset, limit int) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	query := r.readDB.Where("user_id = ?", userID)
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
	query := r.readDB.Model(&model.RecurringTransaction{}).Where("user_id = ?", userID)
	if active != nil {
		query = query.Where("is_active = ?", *active)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RecurringTransactionRepository) Update(rt *model.RecurringTransaction) error {
	return r.writeDB.Save(rt).Error
}

func (r *RecurringTransactionRepository) UpdateWithDB(dbTx *gorm.DB, rt *model.RecurringTransaction) error {
	return dbTx.Save(rt).Error
}

func (r *RecurringTransactionRepository) Delete(id, userID uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Transaction{}).Where("recurring_id = ? AND user_id = ?", id, userID).Update("recurring_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Where("recurring_transaction_id = ?", id).Delete(&model.RecurringTransactionLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecurringTransaction{}).Error
	})
}

func (r *RecurringTransactionRepository) GetDueRecurringTransactions(userID uint64) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := gorm.Expr("NOW()")
	if err := r.readDB.Where("user_id = ? AND is_active = ? AND next_occurrence <= ?", userID, true, now).
		Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

func (r *RecurringTransactionRepository) CreateLog(log *model.RecurringTransactionLog) error {
	return r.writeDB.Create(log).Error
}

func (r *RecurringTransactionRepository) CreateLogWithDB(dbTx *gorm.DB, log *model.RecurringTransactionLog) error {
	return dbTx.Create(log).Error
}

func (r *RecurringTransactionRepository) ListLogs(recurringTransactionID uint64) ([]model.RecurringTransactionLog, error) {
	var logs []model.RecurringTransactionLog
	if err := r.readDB.Where("recurring_transaction_id = ?", recurringTransactionID).
		Order("occurrence_date DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *RecurringTransactionRepository) GetUpcoming(userID uint64, days int) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := time.Now()
	endDate := now.AddDate(0, 0, days)
	if err := r.readDB.Where(
		"user_id = ? AND is_active = ? AND reminder_days > 0 AND next_occurrence <= ? AND date(next_occurrence, '-' || reminder_days || ' days') <= ?",
		userID, true, endDate, endDate.Format("2006-01-02")).
		Order("next_occurrence ASC").Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

func (r *RecurringTransactionRepository) GetAllDue() ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := time.Now()
	if err := r.readDB.Where("is_active = ? AND next_occurrence <= ?", true, now).
		Order("next_occurrence ASC").Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}
