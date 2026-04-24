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

func (r *TransactionLinkRepository) Create(link *model.TransactionJournalLink) error {
	return r.db.Create(link).Error
}

func (r *TransactionLinkRepository) GetByID(id uint64) (*model.TransactionJournalLink, error) {
	var link model.TransactionJournalLink
	if err := r.db.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *TransactionLinkRepository) List(transactionID *uint64, offset, limit int) ([]model.TransactionJournalLink, error) {
	var links []model.TransactionJournalLink
	query := r.db.Model(&model.TransactionJournalLink{})
	if transactionID != nil {
		query = query.Where("transaction_id = ? OR linked_journal_id = ?", *transactionID, *transactionID)
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *TransactionLinkRepository) Count(transactionID *uint64) (int64, error) {
	var count int64
	query := r.db.Model(&model.TransactionJournalLink{})
	if transactionID != nil {
		query = query.Where("transaction_id = ? OR linked_journal_id = ?", *transactionID, *transactionID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransactionLinkRepository) Delete(id uint64) error {
	return r.db.Delete(&model.TransactionJournalLink{}, id).Error
}
