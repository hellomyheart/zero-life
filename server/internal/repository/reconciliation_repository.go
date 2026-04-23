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

func (r *ReconciliationRepository) Create(recon *model.Reconciliation) error {
	return r.db.Create(recon).Error
}

func (r *ReconciliationRepository) GetByAccount(userID, accountID uint64, offset, limit int) ([]model.Reconciliation, error) {
	var recs []model.Reconciliation
	if err := r.db.Where("user_id = ? AND account_id = ?", userID, accountID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *ReconciliationRepository) CountByAccount(userID, accountID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Reconciliation{}).
		Where("user_id = ? AND account_id = ?", userID, accountID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
