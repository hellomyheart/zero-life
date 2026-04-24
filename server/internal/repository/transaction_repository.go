package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// TransactionFilter 交易列表的基础过滤器，用于 List 和 Count 方法。
// 各字段为空或 nil 时表示不过滤该条件。
type TransactionFilter struct {
	Type       string  // 交易类型过滤，如 "withdrawal"（支出）、"deposit"（收入）、"transfer"（转账）
	StartDate  string  // 开始日期过滤，格式 "2006-01-02"
	EndDate    string  // 结束日期过滤，格式 "2006-01-02"
	AccountID  *uint64 // 账户 ID 过滤，匹配源账户或目标账户
	CategoryID *uint64 // 分类 ID 过滤
	TagID      *uint64 // 标签 ID 过滤，通过 JOIN transaction_tags 表实现
}

// TransactionRepository 交易仓库，负责交易记录的数据访问。
// 交易（Transaction）是系统的核心实体，记录每一笔收入、支出和转账。
// 交易可以关联标签（多对多）、拆分为子交易（parent_id）、预加载关联数据。
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 创建交易仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create 创建交易及其关联的标签。使用数据库事务确保原子性。
// 先创建交易记录，再批量创建交易-标签关联记录。
// 执行 SQL（事务内）:
//   1. INSERT INTO transactions (...)
//   2. INSERT INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?), (?, ?), ...
// 参数 txn: 要创建的交易对象，创建后 GORM 会自动填充 ID。
// 参数 tagIDs: 要关联的标签 ID 列表，可以为空。
// 返回: 创建失败时返回错误，事务回滚。
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

// GetByID 根据 ID 和用户 ID 获取单条交易，并预加载所有关联数据。
// 预加载包括：源账户、目标账户、分类、标签、拆分子交易及其标签和分类。
// Preload 会在主查询之后自动发起额外的 SELECT 查询来加载关联数据。
// 执行 SQL:
//   主查询: SELECT * FROM transactions WHERE id = ? AND user_id = ? LIMIT 1
//   预加载: SELECT * FROM accounts WHERE id IN (...)  (Source, Destination)
//           SELECT * FROM categories WHERE id IN (...)  (Category)
//           SELECT * FROM tags JOIN transaction_tags ...  (Tags)
//           SELECT * FROM transactions WHERE parent_id = ?  (Splits)
//           ... 以及 Splits 的 Category 和 Tags
// 参数 id: 交易 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 包含完整关联数据的交易对象。
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

// List 分页获取指定用户的交易列表（仅顶级交易，不含拆分子交易）。
// parent_id IS NULL 过滤条件确保只返回顶级交易，拆分子交易通过 Preload 加载。
// 执行 SQL: SELECT * FROM transactions WHERE user_id = ? AND parent_id IS NULL [AND 过滤条件] ORDER BY date DESC, id DESC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 filter: 交易过滤器，各字段为空时不过滤。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 交易列表（含预加载的源账户、目标账户、分类、标签）。
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

// Count 统计指定用户符合过滤条件的交易总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM transactions WHERE user_id = ? AND parent_id IS NULL [AND 过滤条件]
func (r *TransactionRepository) Count(userID uint64, filter TransactionFilter) (int64, error) {
	var count int64
	query := r.db.Model(&model.Transaction{}).Where("user_id = ? AND parent_id IS NULL", userID)

	query = r.applyFilter(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新交易记录（不含标签）。
// 执行 SQL: UPDATE transactions SET ... WHERE id = ?
// 参数 txn: 要更新的交易对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *TransactionRepository) Update(txn *model.Transaction) error {
	return r.db.Save(txn).Error
}

// UpdateWithTags 更新交易记录及其关联标签。使用数据库事务确保原子性。
// 策略：先更新交易本身，再删除旧的标签关联，最后插入新的标签关联（全量替换）。
// 执行 SQL（事务内）:
//   1. UPDATE transactions SET ... WHERE id = ?
//   2. DELETE FROM transaction_tags WHERE transaction_id = ?
//   3. INSERT INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?), ...
// 参数 txn: 要更新的交易对象。
// 参数 tagIDs: 新的标签 ID 列表，替换原有标签。
// 返回: 更新失败时返回错误，事务回滚。
func (r *TransactionRepository) UpdateWithTags(txn *model.Transaction, tagIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(txn).Error; err != nil {
			return err
		}
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

// Delete 删除交易及其关联的标签和拆分子交易。使用数据库事务确保原子性。
// 删除顺序：1.交易标签关联 → 2.拆分子交易 → 3.交易本身
// 执行 SQL（事务内）:
//   1. DELETE FROM transaction_tags WHERE transaction_id = ?
//   2. DELETE FROM transactions WHERE parent_id = ?
//   3. DELETE FROM transactions WHERE id = ? AND user_id = ?
// 参数 id: 交易 ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的交易。
// 返回: 删除失败时返回错误，事务回滚。
func (r *TransactionRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("transaction_id = ?", id).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		if err := tx.Where("parent_id = ?", id).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// Search 根据关键词搜索交易，匹配描述（description）或备注（notes）字段。
// 使用 LIKE 模糊匹配，仅返回顶级交易（不含拆分子交易）。
// 执行 SQL: SELECT * FROM transactions WHERE user_id = ? AND parent_id IS NULL AND (description LIKE ? OR notes LIKE ?) ORDER BY date DESC, id DESC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 keyword: 搜索关键词。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 匹配的交易列表（含预加载的关联数据）。
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

// SearchCount 统计关键词搜索结果的交易总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM transactions WHERE user_id = ? AND parent_id IS NULL AND (description LIKE ? OR notes LIKE ?)
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

// GetForAudit 获取指定账户在时间范围内的交易，用于对账审计。
// 匹配源账户或目标账户，可选按日期范围和是否已对账过滤。
// 执行 SQL: SELECT * FROM transactions WHERE source_id = ? OR destination_id = ? [AND date >= ?] [AND date <= ?] [AND is_reconciled = ?] ORDER BY date ASC, id ASC
// 参数 accountID: 账户 ID。
// 参数 startDate: 开始日期字符串，格式 "2006-01-02"，为空不过滤。
// 参数 endDate: 结束日期字符串，格式 "2006-01-02"，为空不过滤。
// 参数 reconciled: 是否已对账过滤，为 nil 不过滤。
// 返回: 符合条件的交易列表（含预加载的关联数据）。
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

// applyFilter 将 TransactionFilter 中的过滤条件应用到 GORM 查询上。
// 这是一个内部辅助方法，根据过滤器的非空字段动态添加 WHERE 条件。
// 标签过滤通过 JOIN transaction_tags 关联表实现。
// 参数 query: 要添加条件的 GORM 查询对象。
// 参数 filter: 过滤器。
// 返回: 添加了过滤条件后的查询对象。
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
	if err := r.db.Where("transaction_id = ?", txnID).Delete(&model.TransactionTag{}).Error; err != nil {
		return err
	}

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

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("description LIKE ? OR notes LIKE ?", keyword, keyword)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.StartDate != "" {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("date <= ?", filter.EndDate)
	}

	if filter.MinAmount != "" {
		query = query.Where("amount >= ?", filter.MinAmount)
	}
	if filter.MaxAmount != "" {
		query = query.Where("amount <= ?", filter.MaxAmount)
	}

	if filter.AccountID != nil {
		query = query.Where("source_id = ? OR destination_id = ?", *filter.AccountID, *filter.AccountID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id").
			Where("transaction_tags.tag_id = ?", *filter.TagID)
	}

	orderClause := "date DESC"
	if filter.Sort != "" {
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

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("description LIKE ? OR notes LIKE ?", keyword, keyword)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.StartDate != "" {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("date <= ?", filter.EndDate)
	}

	if filter.MinAmount != "" {
		query = query.Where("amount >= ?", filter.MinAmount)
	}
	if filter.MaxAmount != "" {
		query = query.Where("amount <= ?", filter.MaxAmount)
	}

	if filter.AccountID != nil {
		query = query.Where("source_id = ? OR destination_id = ?", *filter.AccountID, *filter.AccountID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id").
			Where("transaction_tags.tag_id = ?", *filter.TagID)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}
