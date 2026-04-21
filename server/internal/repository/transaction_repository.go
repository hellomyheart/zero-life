package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type TransactionFilter struct {
	Type       string
	StartDate  string
	EndDate    string
	AccountID  *uint64
	CategoryID *uint64
	TagID      *uint64
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(txn *model.Transaction, tagIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(txn).Error; err != nil {
			return err
		}
		if len(tagIDs) > 0 {
			transactionTags := make([]model.TransactionTag, 0, len(tagIDs))
			for _, tagID := range tagIDs {
				transactionTags = append(transactionTags, model.TransactionTag{
					TransactionID: txn.ID,
					TagID:         tagID,
				})
			}
			if err := tx.Create(&transactionTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TransactionRepository) GetByID(id, userID uint64) (*model.Transaction, error) {
	var txn model.Transaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Source").
		Preload("Destination").
		Preload("Category").
		Preload("Tags").
		Preload("Splits").
		Preload("Splits.Tags").
		Preload("Splits.Category").
		First(&txn).Error; err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *TransactionRepository) List(userID uint64, filter TransactionFilter, offset, limit int) ([]model.Transaction, error) {
	var txns []model.Transaction
	query := r.db.Where("user_id = ? AND parent_id IS NULL", userID)

	query = r.applyFilter(query, filter)

	if err := query.Preload("Source").Preload("Destination").Preload("Category").Preload("Tags").
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
}

func (r *TransactionRepository) Count(userID uint64, filter TransactionFilter) (int64, error) {
	var count int64
	query := r.db.Model(&model.Transaction{}).Where("user_id = ? AND parent_id IS NULL", userID)

	query = r.applyFilter(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransactionRepository) Update(txn *model.Transaction) error {
	return r.db.Save(txn).Error
}

func (r *TransactionRepository) UpdateWithTags(txn *model.Transaction, tagIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(txn).Error; err != nil {
			return err
		}
		// Replace tags: delete old, insert new
		if err := tx.Where("transaction_id = ?", txn.ID).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		if len(tagIDs) > 0 {
			transactionTags := make([]model.TransactionTag, 0, len(tagIDs))
			for _, tagID := range tagIDs {
				transactionTags = append(transactionTags, model.TransactionTag{
					TransactionID: txn.ID,
					TagID:         tagID,
				})
			}
			if err := tx.Create(&transactionTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TransactionRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete transaction tags
		if err := tx.Where("transaction_id = ?", id).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		// Delete splits (child transactions)
		if err := tx.Where("parent_id = ?", id).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		// Delete the transaction itself
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TransactionRepository) Search(userID uint64, keyword string, offset, limit int) ([]model.Transaction, error) {
	var txns []model.Transaction
	like := "%" + keyword + "%"
	if err := r.db.Where("user_id = ? AND parent_id IS NULL", userID).
		Where("description LIKE ? OR notes LIKE ?", like, like).
		Preload("Source").Preload("Destination").Preload("Category").Preload("Tags").
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
}

func (r *TransactionRepository) SearchCount(userID uint64, keyword string) (int64, error) {
	var count int64
	like := "%" + keyword + "%"
	if err := r.db.Model(&model.Transaction{}).
		Where("user_id = ? AND parent_id IS NULL", userID).
		Where("description LIKE ? OR notes LIKE ?", like, like).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransactionRepository) applyFilter(query *gorm.DB, filter TransactionFilter) *gorm.DB {
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.StartDate != "" {
		if t, err := time.Parse("2006-01-02", filter.StartDate); err == nil {
			query = query.Where("date >= ?", t)
		}
	}
	if filter.EndDate != "" {
		if t, err := time.Parse("2006-01-02", filter.EndDate); err == nil {
			query = query.Where("date <= ?", t)
		}
	}
	if filter.AccountID != nil {
		query = query.Where("source_id = ? OR destination_id = ?", *filter.AccountID, *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id AND transaction_tags.tag_id = ?", *filter.TagID)
	}
	return query
}
