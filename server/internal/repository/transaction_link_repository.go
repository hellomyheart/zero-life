package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type TransactionLinkRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewTransactionLinkRepository(readDB, writeDB *gorm.DB) *TransactionLinkRepository {
	return &TransactionLinkRepository{readDB: readDB, writeDB: writeDB}
}

func (r *TransactionLinkRepository) Create(link *model.TransactionJournalLink) error {
	return r.writeDB.Create(link).Error
}

func (r *TransactionLinkRepository) GetByID(id, userID uint64) (*model.TransactionJournalLink, error) {
	var link model.TransactionJournalLink
	if err := r.readDB.Joins("JOIN transactions ON transactions.id = transaction_journal_links.transaction_id").
		Where("transaction_journal_links.id = ? AND transactions.user_id = ?", id, userID).
		First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *TransactionLinkRepository) List(userID uint64, transactionID *uint64, offset, limit int) ([]model.TransactionJournalLink, error) {
	var links []model.TransactionJournalLink
	query := r.readDB.Model(&model.TransactionJournalLink{}).
		Joins("JOIN transactions ON transactions.id = transaction_journal_links.transaction_id").
		Where("transactions.user_id = ?", userID)
	if transactionID != nil {
		query = query.Where("transaction_journal_links.transaction_id = ? OR transaction_journal_links.linked_journal_id = ?", *transactionID, *transactionID)
	}
	if err := query.Order("transaction_journal_links.created_at DESC").Offset(offset).Limit(limit).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *TransactionLinkRepository) Count(userID uint64, transactionID *uint64) (int64, error) {
	var count int64
	query := r.readDB.Model(&model.TransactionJournalLink{}).
		Joins("JOIN transactions ON transactions.id = transaction_journal_links.transaction_id").
		Where("transactions.user_id = ?", userID)
	if transactionID != nil {
		query = query.Where("transaction_journal_links.transaction_id = ? OR transaction_journal_links.linked_journal_id = ?", *transactionID, *transactionID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransactionLinkRepository) Delete(id, userID uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.TransactionJournalLink{}).
			Joins("JOIN transactions ON transactions.id = transaction_journal_links.transaction_id").
			Where("transaction_journal_links.id = ? AND transactions.user_id = ?", id, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Delete(&model.TransactionJournalLink{}, id).Error
	})
}
