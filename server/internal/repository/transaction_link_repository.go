// Package repository 提供数据访问层，封装所有与数据库的交互逻辑。
// 每个 Repository 结构体对应一个数据库表，提供增删改查等基本操作。
// 所有 Repository 都依赖 GORM 作为 ORM 框架，通过 *gorm.DB 执行数据库操作。
package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// TransactionLinkRepository 交易链接仓库，负责交易之间关联关系的数据访问。
// 交易链接（TransactionJournalLink）用于表示两笔交易之间的关联，例如"转账-收款"配对。
type TransactionLinkRepository struct {
	db *gorm.DB
}

// NewTransactionLinkRepository 创建交易链接仓库实例。
// 参数 db: GORM 数据库连接实例。
// 返回: 初始化后的 TransactionLinkRepository 指针。
func NewTransactionLinkRepository(db *gorm.DB) *TransactionLinkRepository {
	return &TransactionLinkRepository{db: db}
}

// Create 创建一条新的交易链接记录。
// 执行 SQL: INSERT INTO transaction_journal_links (...)
// 参数 link: 要创建的交易链接对象，GORM 会自动填充 ID、CreatedAt 等字段。
// 返回: 创建失败时返回错误。
func (r *TransactionLinkRepository) Create(link *model.TransactionJournalLink) error {
	return r.db.Create(link).Error
}

// GetByID 根据 ID 和用户 ID 获取单条交易链接。
// 通过 JOIN transactions 表来验证该链接属于指定用户（权限校验）。
// 执行 SQL: SELECT * FROM transaction_journal_links JOIN transactions ON ... WHERE transaction_journal_links.id = ? AND transactions.user_id = ? LIMIT 1
// 参数 id: 交易链接 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验，确保用户只能查看自己的数据。
// 返回: 找到的交易链接对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *TransactionLinkRepository) GetByID(id, userID uint64) (*model.TransactionJournalLink, error) {
	var link model.TransactionJournalLink
	if err := r.db.Joins("JOIN transactions ON transactions.id = transaction_journal_links.transaction_id").
		Where("transaction_journal_links.id = ? AND transactions.user_id = ?", id, userID).
		First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

// List 分页获取指定用户的交易链接列表。
// 通过 JOIN transactions 表过滤出属于该用户的链接，可选按交易 ID 过滤。
// 执行 SQL: SELECT * FROM transaction_journal_links JOIN transactions ON ... WHERE transactions.user_id = ? [AND transaction_id = ? OR linked_journal_id = ?] ORDER BY created_at DESC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 transactionID: 可选的交易 ID 过滤条件，为 nil 时不按交易过滤。
// 参数 offset: 分页偏移量（跳过的记录数）。
// 参数 limit: 每页记录数。
// 返回: 交易链接列表。
func (r *TransactionLinkRepository) List(userID uint64, transactionID *uint64, offset, limit int) ([]model.TransactionJournalLink, error) {
	var links []model.TransactionJournalLink
	query := r.db.Model(&model.TransactionJournalLink{}).
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

// Count 统计指定用户的交易链接总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM transaction_journal_links JOIN transactions ON ... WHERE transactions.user_id = ? [AND ...]
// 参数与 List 相同，返回符合条件的记录总数。
func (r *TransactionLinkRepository) Count(userID uint64, transactionID *uint64) (int64, error) {
	var count int64
	query := r.db.Model(&model.TransactionJournalLink{}).
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

// Delete 根据 ID 删除一条交易链接记录。
// 注意：此方法不做用户权限校验，仅按 ID 删除，调用方需自行确保权限。
// 执行 SQL: DELETE FROM transaction_journal_links WHERE id = ?
// 参数 id: 要删除的交易链接 ID。
// 返回: 删除失败时返回错误。
func (r *TransactionLinkRepository) Delete(id uint64) error {
	return r.db.Delete(&model.TransactionJournalLink{}, id).Error
}
