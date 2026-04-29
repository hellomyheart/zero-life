package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// TransactionFilter 交易列表的基础过滤器，用于 List 和 Count 方法。
// 各字段为空或 nil 时表示不过滤该条件。
type TransactionFilter struct {
	Type        string   // 交易类型过滤，如 "withdrawal"（支出）、"deposit"（收入）、"transfer"（转账）
	StartDate   string   // 开始日期过滤，格式 "2006-01-02"
	EndDate     string   // 结束日期过滤，格式 "2006-01-02"
	AccountID   *uint64  // 账户 ID 过滤，匹配源账户或目标账户
	CategoryID  *uint64  // 分类 ID 过滤（单个）
	CategoryIDs []uint64 // 分类 ID 过滤（多个，OR 关系）
	TagID       *uint64  // 标签 ID 过滤（单个），通过 JOIN transaction_tags 表实现
	TagIDs      []uint64 // 标签 ID 过滤（多个，OR 关系），通过 JOIN transaction_tags 表实现
	Keyword     string   // 关键词搜索，匹配描述或备注
}

// TransactionRepository 交易仓库，负责交易记录的数据访问。
// 交易（Transaction）是系统的核心实体，记录每一笔收入、支出和转账。
// 交易可以关联标签（多对多）、拆分为子交易（parent_id）、预加载关联数据。
type TransactionRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewTransactionRepository 创建交易仓库实例。
// 参数 readDB: 读库（只读连接池，并发安全）
// 参数 writeDB: 写库（单连接，串行保证安全）
func NewTransactionRepository(readDB, writeDB *gorm.DB) *TransactionRepository {
	return &TransactionRepository{readDB: readDB, writeDB: writeDB}
}

// Create 创建交易及其关联的标签。
// 支持传入事务对象，确保与余额更新在同一事务中执行。
// 先创建交易记录，再批量创建交易-标签关联记录。
// 参数 txn: 要创建的交易对象，创建后 GORM 会自动填充 ID。
// 参数 tagIDs: 要关联的标签 ID 列表，可以为空。
// 返回: 创建失败时返回错误。
func (r *TransactionRepository) Create(txn *model.Transaction, tagIDs []uint64) error {
	return r.CreateWithDB(r.writeDB, txn, tagIDs)
}

// CreateWithDB 使用指定的 DB 对象创建交易及其关联的标签。
// 用于在事务中创建交易，确保与余额更新在同一事务中执行。
func (r *TransactionRepository) CreateWithDB(db *gorm.DB, txn *model.Transaction, tagIDs []uint64) error {
	if err := db.Create(txn).Error; err != nil {
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
		if err := db.Create(&transactionTags).Error; err != nil {
			return err
		}
	}
	return nil
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
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).
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
	query := r.readDB.Where("user_id = ? AND parent_id IS NULL", userID)

	query = r.applyFilter(query, filter)

	if err := query.Preload("Source").Preload("Destination").Preload("Category").Preload("Tags").
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
}

// ListAll 查询指定用户符合过滤条件的所有交易（无条数限制）。
// 内部循环分页查询，每次取 5000 条，直到取完所有数据。
// 用于预算支出计算、报表统计等需要完整数据的场景。
func (r *TransactionRepository) ListAll(userID uint64, filter TransactionFilter) ([]model.Transaction, error) {
	var allTxns []model.Transaction
	const batchSize = 5000
	offset := 0
	for {
		var batch []model.Transaction
		query := r.readDB.Where("user_id = ? AND parent_id IS NULL", userID)
		query = r.applyFilter(query, filter)
		if err := query.Preload("Source").Preload("Destination").Preload("Category").Preload("Tags").
			Order("date DESC, id DESC").Offset(offset).Limit(batchSize).Find(&batch).Error; err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		allTxns = append(allTxns, batch...)
		if len(batch) < batchSize {
			break
		}
		offset += batchSize
	}
	return allTxns, nil
}

// Count 统计指定用户符合过滤条件的交易总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM transactions WHERE user_id = ? AND parent_id IS NULL [AND 过滤条件]
func (r *TransactionRepository) Count(userID uint64, filter TransactionFilter) (int64, error) {
	var count int64
	query := r.readDB.Model(&model.Transaction{}).Where("user_id = ? AND parent_id IS NULL", userID)

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
	return r.writeDB.Save(txn).Error
}

// UpdateWithTags 更新交易记录及其关联标签。
// 策略：先更新交易本身，再删除旧的标签关联，最后插入新的标签关联（全量替换）。
// 参数 txn: 要更新的交易对象。
// 参数 tagIDs: 新的标签 ID 列表，替换原有标签。
// 返回: 更新失败时返回错误。
func (r *TransactionRepository) UpdateWithTags(txn *model.Transaction, tagIDs []uint64) error {
	return r.UpdateWithTagsAndDB(r.writeDB, txn, tagIDs)
}

// UpdateWithTagsAndDB 使用指定的 DB 对象更新交易记录及其关联标签。
// 用于在事务中更新交易，确保与余额更新在同一事务中执行。
func (r *TransactionRepository) UpdateWithTagsAndDB(db *gorm.DB, txn *model.Transaction, tagIDs []uint64) error {
	if err := db.Save(txn).Error; err != nil {
		return err
	}
	if err := db.Where("transaction_id = ?", txn.ID).Delete(&model.TransactionTag{}).Error; err != nil {
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
		if err := db.Create(&transactionTags).Error; err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除交易及其关联的标签和拆分子交易。
// 使用硬删除（Unscoped），因为余额已硬回滚，软删除记录会导致恢复时余额不一致。
func (r *TransactionRepository) Delete(id, userID uint64) error {
	return r.DeleteWithDB(r.writeDB, id, userID)
}

// DeleteWithDB 使用指定的 DB 对象删除交易及其关联的标签和拆分子交易。
// 使用硬删除（Unscoped），因为余额已在同一事务中硬回滚，
// 软删除记录留在表中会导致 Unscoped 恢复时余额不一致。
func (r *TransactionRepository) DeleteWithDB(db *gorm.DB, id, userID uint64) error {
	if err := db.Where("transaction_id = ?", id).Delete(&model.TransactionTag{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("parent_id = ?", id).Delete(&model.Transaction{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&model.Transaction{}).Error; err != nil {
		return err
	}
	return nil
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
	if err := r.readDB.Where("user_id = ? AND parent_id IS NULL", userID).
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
	if err := r.readDB.Model(&model.Transaction{}).
		Where("user_id = ? AND parent_id IS NULL", userID).
		Where("description LIKE ? OR notes LIKE ?", like, like).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetForAudit 获取指定账户在时间范围内的交易，用于对账审计。
// 匹配源账户或目标账户，可选按日期范围和是否已对账过滤。
// 修复：添加 user_id 过滤条件，防止跨用户数据泄露。
// 执行 SQL: SELECT * FROM transactions WHERE user_id = ? AND (source_id = ? OR destination_id = ?) [AND date >= ?] [AND date <= ?] [AND is_reconciled = ?] ORDER BY date ASC, id ASC
// 参数 userID: 用户 ID，用于权限校验，确保只能查询自己的交易。
// 参数 accountID: 账户 ID。
// 参数 startDate: 开始日期字符串，格式 "2006-01-02"，为空不过滤。
// 参数 endDate: 结束日期字符串，格式 "2006-01-02"，为空不过滤。
// 参数 reconciled: 是否已对账过滤，为 nil 不过滤。
// 返回: 符合条件的交易列表（含预加载的关联数据）。
func (r *TransactionRepository) GetForAudit(userID, accountID uint64, startDate, endDate string, reconciled *bool) ([]model.Transaction, error) {
	var txns []model.Transaction
	query := r.readDB.Where("user_id = ? AND (source_id = ? OR destination_id = ?)", userID, accountID, accountID)

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
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("description LIKE ? OR notes LIKE ?", keyword, keyword)
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
	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}
	if filter.TagID != nil && len(filter.TagIDs) > 0 {
		allTagIDs := make([]uint64, 0, 1+len(filter.TagIDs))
		allTagIDs = append(allTagIDs, *filter.TagID)
		allTagIDs = append(allTagIDs, filter.TagIDs...)
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id AND transaction_tags.tag_id IN ?", allTagIDs)
	} else if filter.TagID != nil {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id AND transaction_tags.tag_id = ?", *filter.TagID)
	} else if len(filter.TagIDs) > 0 {
		query = query.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id AND transaction_tags.tag_id IN ?", filter.TagIDs)
	}
	return query
}

// GetByAccountAndDateRange 获取指定账户在时间范围内的交易
func (r *TransactionRepository) GetByAccountAndDateRange(userID, accountID uint64, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.readDB.Where("user_id = ? AND (source_id = ? OR destination_id = ?) AND date >= ? AND date <= ?",
		userID, accountID, accountID, startDate, endDate).
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetByDateRange 获取指定时间范围内的所有交易
func (r *TransactionRepository) GetByDateRange(userID uint64, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.readDB.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Preload("Tags").
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetByTypeAndDateRange 获取指定类型和时间范围内的交易
func (r *TransactionRepository) GetByTypeAndDateRange(userID uint64, txnType string, startDate, endDate time.Time) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.readDB.Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
		userID, txnType, startDate, endDate).
		Preload("Tags").
		Order("date asc").
		Find(&txns).Error
	return txns, err
}

// GetSplits 获取拆分交易列表
func (r *TransactionRepository) GetSplits(parentID, userID uint64) ([]model.Transaction, error) {
	var splits []model.Transaction
	err := r.readDB.Where("parent_id = ? AND user_id = ?", parentID, userID).
		Preload("Category").
		Preload("Tags").
		Order("id asc").
		Find(&splits).Error
	return splits, err
}

// CreateBatch 批量创建交易
func (r *TransactionRepository) CreateBatch(txns []model.Transaction) error {
	return r.CreateBatchWithDB(r.writeDB, txns)
}

// CreateBatchWithDB 使用指定的 DB 对象批量创建交易
func (r *TransactionRepository) CreateBatchWithDB(db *gorm.DB, txns []model.Transaction) error {
	return db.Create(&txns).Error
}

// DeleteBatch 批量删除交易及其关联的标签和拆分子交易。使用数据库事务确保原子性。
// 使用硬删除（Unscoped），因为余额已在同一事务中硬回滚，
// 软删除记录留在表中会导致 Unscoped 恢复时余额不一致。
func (r *TransactionRepository) DeleteBatch(ids []uint64, userID uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("transaction_id IN ?", ids).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("parent_id IN ?", ids).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.Transaction{}).Error
	})
}

// AttachTags 为交易添加标签（全量替换：先删除旧标签，再插入新标签）
func (r *TransactionRepository) AttachTags(txnID uint64, tagIDs []uint64) error {
	return r.AttachTagsWithDB(r.writeDB, txnID, tagIDs)
}

// AttachTagsWithDB 使用指定的 DB 对象为交易添加标签（全量替换）
func (r *TransactionRepository) AttachTagsWithDB(db *gorm.DB, txnID uint64, tagIDs []uint64) error {
	if err := db.Where("transaction_id = ?", txnID).Delete(&model.TransactionTag{}).Error; err != nil {
		return err
	}

	transactionTags := make([]model.TransactionTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		transactionTags = append(transactionTags, model.TransactionTag{
			TransactionID: txnID,
			TagID:         tagID,
		})
	}

	return db.Create(&transactionTags).Error
}

// AddTags 为交易追加标签（仅添加不存在的标签，不删除已有标签）
func (r *TransactionRepository) AddTags(txnID uint64, tagIDs []uint64) error {
	return r.AddTagsWithDB(r.writeDB, txnID, tagIDs)
}

// AddTagsWithDB 使用指定的 DB 对象为交易追加标签
func (r *TransactionRepository) AddTagsWithDB(db *gorm.DB, txnID uint64, tagIDs []uint64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	var existingIDs []uint64
	db.Model(&model.TransactionTag{}).Where("transaction_id = ? AND tag_id IN ?", txnID, tagIDs).
		Pluck("tag_id", &existingIDs)
	existingSet := make(map[uint64]bool, len(existingIDs))
	for _, id := range existingIDs {
		existingSet[id] = true
	}
	newTags := make([]model.TransactionTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		if !existingSet[tagID] {
			newTags = append(newTags, model.TransactionTag{
				TransactionID: txnID,
				TagID:         tagID,
			})
		}
	}
	if len(newTags) == 0 {
		return nil
	}
	return db.Create(&newTags).Error
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
// 修复：添加 parent_id IS NULL 过滤条件，避免返回拆分子交易（与 List 方法保持一致）。
func (r *TransactionRepository) AdvancedSearch(userID uint64, filter AdvancedSearchFilter, offset, limit int) ([]model.Transaction, error) {
	query := r.readDB.Model(&model.Transaction{}).Where("user_id = ? AND parent_id IS NULL", userID)

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
// 修复：添加 parent_id IS NULL 过滤条件，与 AdvancedSearch 保持一致。
func (r *TransactionRepository) AdvancedSearchCount(userID uint64, filter AdvancedSearchFilter) (int64, error) {
	query := r.readDB.Model(&model.Transaction{}).Where("user_id = ? AND parent_id IS NULL", userID)

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
