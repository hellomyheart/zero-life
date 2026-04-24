package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type ReconciliationRepository struct {
	db *gorm.DB
}

func NewReconciliationRepository(db *gorm.DB) *ReconciliationRepository {
	return &ReconciliationRepository{db: db}
}

func (r *ReconciliationRepository) Create(rec *model.TransactionReconciliation) error {
	return r.db.Create(rec).Error
}

func (r *ReconciliationRepository) GetByID(id uint64) (*model.TransactionReconciliation, error) {
	var rec model.TransactionReconciliation
	if err := r.db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *ReconciliationRepository) List(accountID *uint64, offset, limit int) ([]model.TransactionReconciliation, error) {
	var recs []model.TransactionReconciliation
	query := r.db.Model(&model.TransactionReconciliation{})
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *ReconciliationRepository) Count(accountID *uint64) (int64, error) {
	var count int64
	query := r.db.Model(&model.TransactionReconciliation{})
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ReconciliationRepository) Update(rec *model.TransactionReconciliation) error {
	return r.db.Save(rec).Error
}

func (r *ReconciliationRepository) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("reconciliation_id = ?", id).Delete(&model.ReconciliationEntry{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.TransactionReconciliation{}, id).Error
	})
}

func (r *ReconciliationRepository) CreateEntry(entry *model.ReconciliationEntry) error {
	return r.db.Create(entry).Error
}

func (r *ReconciliationRepository) ListEntries(reconciliationID uint64) ([]model.ReconciliationEntry, error) {
	var entries []model.ReconciliationEntry
	if err := r.db.Where("reconciliation_id = ?", reconciliationID).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}
