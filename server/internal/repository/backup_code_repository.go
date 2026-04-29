package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type BackupCodeRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBackupCodeRepository(readDB, writeDB *gorm.DB) *BackupCodeRepository {
	return &BackupCodeRepository{readDB: readDB, writeDB: writeDB}
}

func (r *BackupCodeRepository) Create(codes []model.BackupCode) error {
	return r.writeDB.Create(&codes).Error
}

func (r *BackupCodeRepository) ListByUserID(userID uint64) ([]model.BackupCode, error) {
	var codes []model.BackupCode
	if err := r.readDB.Where("user_id = ?", userID).Order("id ASC").Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *BackupCodeRepository) FindByCode(userID uint64, code string) (*model.BackupCode, error) {
	var bc model.BackupCode
	if err := r.readDB.Where("user_id = ? AND code = ? AND used_at IS NULL", userID, code).First(&bc).Error; err != nil {
		return nil, err
	}
	return &bc, nil
}

func (r *BackupCodeRepository) MarkUsed(id uint64) error {
	now := time.Now()
	return r.writeDB.Model(&model.BackupCode{}).Where("id = ?", id).Update("used_at", now).Error
}

func (r *BackupCodeRepository) DeleteByUserID(userID uint64) error {
	return r.writeDB.Where("user_id = ?", userID).Delete(&model.BackupCode{}).Error
}
