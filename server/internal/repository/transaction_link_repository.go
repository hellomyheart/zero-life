package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type TransactionLinkRepository struct {
	db *gorm.DB
}

func NewTransactionLinkRepository(db *gorm.DB) *TransactionLinkRepository {
	return &TransactionLinkRepository{db: db}
}

func (r *TransactionLinkRepository) Create(link *model.TransactionLink) error {
	return r.db.Create(link).Error
}

func (r *TransactionLinkRepository) ListByTransactionID(transactionID uint64) ([]model.TransactionLink, error) {
	var links []model.TransactionLink
	if err := r.db.Where("source_id = ? OR destination_id = ?", transactionID, transactionID).
		Preload("LinkType").
		Order("id ASC").
		Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *TransactionLinkRepository) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&model.TransactionLink{}).Error
}

func (r *TransactionLinkRepository) Exists(linkTypeID, sourceID, destinationID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.TransactionLink{}).
		Where("link_type_id = ? AND source_id = ? AND destination_id = ?", linkTypeID, sourceID, destinationID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TransactionLinkRepository) GetByID(id uint64) (*model.TransactionLink, error) {
	var link model.TransactionLink
	if err := r.db.Where("id = ?", id).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}
