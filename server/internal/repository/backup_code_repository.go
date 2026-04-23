package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type BackupCodeRepository struct {
	db *gorm.DB
}

func NewBackupCodeRepository(db *gorm.DB) *BackupCodeRepository {
	return &BackupCodeRepository{db: db}
}

func (r *BackupCodeRepository) Create(codes []model.BackupCode) error {
	return r.db.Create(&codes).Error
}

func (r *BackupCodeRepository) ListByUserID(userID uint64) ([]model.BackupCode, error) {
	var codes []model.BackupCode
	if err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *BackupCodeRepository) FindByCode(userID uint64, code string) (*model.BackupCode, error) {
	var bc model.BackupCode
	if err := r.db.Where("user_id = ? AND code = ? AND used_at IS NULL", userID, code).First(&bc).Error; err != nil {
		return nil, err
	}
	return &bc, nil
}

func (r *BackupCodeRepository) MarkUsed(id uint64) error {
	now := time.Now()
	return r.db.Model(&model.BackupCode{}).Where("id = ?", id).Update("used_at", now).Error
}

func (r *BackupCodeRepository) DeleteByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.BackupCode{}).Error
}
