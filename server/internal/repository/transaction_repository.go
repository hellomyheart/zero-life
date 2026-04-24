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

func (r *TransactionRepository) GetForAudit(accountID uint64, startDate, endDate string, reconciled *bool) ([]model.Transaction, error) {
	var txns []model.Transaction
	query := r.db.Where("source_id = ? OR destination_id = ?", accountID, accountID)

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("date >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("date <= ?", t)
		}
	}
	if reconciled != nil {
		query = query.Where("is_reconciled = ?", *reconciled)
	}

	if err := query.Preload("Source").Preload("Destination").Preload("Category").Preload("Tags").
		Order("date ASC, id ASC").Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
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

// GetByAccountAndDateRange 获取指定账户在时间范围内的交易
func (r *TransactionRepository) GetByAccountAndDateRange(userID, accountID uint64, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Where("user_id = ? AND (source_id = ? OR destination_id = ?) AND date >= ? AND date <= ?",
		userID, accountID, accountID, startDate, endDate).
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetByDateRange 获取指定时间范围内的所有交易
func (r *TransactionRepository) GetByDateRange(userID uint64, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Preload("Tags").
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetByTypeAndDateRange 获取指定类型和时间范围内的交易
func (r *TransactionRepository) GetByTypeAndDateRange(userID uint64, txnType string, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
		userID, txnType, startDate, endDate).
		Preload("Tags").
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetSplits 获取拆分交易列表
func (r *TransactionRepository) GetSplits(parentID, userID uint64) ([]model.Transaction, error) {
	var splits []model.Transaction
	err := r.db.Where("parent_id = ? AND user_id = ?", parentID, userID).
		Preload("Category").
		Preload("Tags").
		Order("id asc").
		Find(&splits).Error
	return splits, err
}

// CreateBatch 批量创建交易
func (r *TransactionRepository) CreateBatch(txns []model.Transaction) error {
	return r.db.Create(&txns).Error
}

// DeleteBatch 批量删除交易
func (r *TransactionRepository) DeleteBatch(ids []uint64, userID uint64) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.Transaction{}).Error
}

// AttachTags 为交易添加标签
func (r *TransactionRepository) AttachTags(txnID uint64, tagIDs []uint64) error {
	// 删除现有标签关联
	if err := r.db.Where("transaction_id = ?", txnID).Delete(&model.TransactionTag{}).Error; err != nil {
		return err
	}

	// 添加新标签关联
	transactionTags := make([]model.TransactionTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		transactionTags = append(transactionTags, model.TransactionTag{
			TransactionID: txnID,
			TagID:         tagID,
		})
	}

	return r.db.Create(&transactionTags).Error
}

// AdvancedSearchFilter 高级搜索过滤器
type AdvancedSearchFilter struct {
	Keyword    string  // 关键词（描述、备注）
	Type       string  // 交易类型
	StartDate  string  // 开始日期
	EndDate    string  // 结束日期
	MinAmount  string  // 最小金额
	MaxAmount  string  // 最大金额
	AccountID  *uint64 // 账户ID
	CategoryID *uint64 // 分类ID
	TagID      *uint64 // 标签ID
	Sort       string  // 排序字段
}

// AdvancedSearch 高级搜索交易
func (r *TransactionRepository) AdvancedSearch(userID uint64, filter AdvancedSearchFilter, offset, limit int) ([]model.Transaction, error) {
	query := r.db.Model(&model.Transaction{}).Where("user_id = ?", userID)

	// 关键词搜索（描述或备注）
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("description LIKE ? OR notes LIKE ?", keyword, keyword)
	}

	// 交易类型过滤
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	// 日期范围过滤
	if filter.StartDate != "" {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("date <= ?", filter.EndDate)
	}

	// 金额范围过滤
	if filter.MinAmount != "" {
		query = query.Where("amount >= ?", filter.MinAmount)
	}
	if filter.MaxAmount != "" {
		query = query.Where("amount <= ?", filter.MaxAmount)
	}

	// 账户过滤
	if filter.AccountID != nil {
		query = query.Where("source_id = ? OR destination_id = ?", *filter.AccountID, *filter.AccountID)
	}

	// 分类过滤
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	// 标签过滤
	if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id").
			Where("transaction_tags.tag_id = ?", *filter.TagID)
	}

	// 排序
	orderClause := "date DESC"
	if filter.Sort != "" {
		// 支持的排序字段: date, -date, amount, -amount, created_at, -created_at
		switch filter.Sort {
		case "date":
			orderClause = "date ASC"
		case "-date":
			orderClause = "date DESC"
		case "amount":
			orderClause = "amount ASC"
		case "-amount":
			orderClause = "amount DESC"
		case "created_at":
			orderClause = "created_at ASC"
		case "-created_at":
			orderClause = "created_at DESC"
		}
	}

	var txns []model.Transaction
	err := query.Preload("Source").
		Preload("Destination").
		Preload("Category").
		Preload("Tags").
		Order(orderClause).
		Offset(offset).
		Limit(limit).
		Find(&txns).Error

	return txns, err
}

// AdvancedSearchCount 高级搜索计数
func (r *TransactionRepository) AdvancedSearchCount(userID uint64, filter AdvancedSearchFilter) (int64, error) {
	query := r.db.Model(&model.Transaction{}).Where("user_id = ?", userID)

	// 关键词搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("description LIKE ? OR notes LIKE ?", keyword, keyword)
	}

	// 交易类型过滤
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	// 日期范围过滤
	if filter.StartDate != "" {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("date <= ?", filter.EndDate)
	}

	// 金额范围过滤
	if filter.MinAmount != "" {
		query = query.Where("amount >= ?", filter.MinAmount)
	}
	if filter.MaxAmount != "" {
		query = query.Where("amount <= ?", filter.MaxAmount)
	}

	// 账户过滤
	if filter.AccountID != nil {
		query = query.Where("source_id = ? OR destination_id = ?", *filter.AccountID, *filter.AccountID)
	}

	// 分类过滤
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	// 标签过滤
	if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id").
			Where("transaction_tags.tag_id = ?", *filter.TagID)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}
