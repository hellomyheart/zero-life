package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// ReconciliationRepository 对账仓库，负责账户对账记录的数据访问。
// 对账（Reconciliation）是将银行账单与系统中的交易记录进行核对的过程，
// 确保账户余额与实际银行余额一致。每条对账记录包含多条对账条目（ReconciliationEntry）。
type ReconciliationRepository struct {
	db *gorm.DB
}

// NewReconciliationRepository 创建对账仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewReconciliationRepository(db *gorm.DB) *ReconciliationRepository {
	return &ReconciliationRepository{db: db}
}

// Create 创建一条新的对账记录。
// 执行 SQL: INSERT INTO transaction_reconciliations (...)
// 参数 rec: 要创建的对账记录对象。
// 返回: 创建失败时返回错误。
func (r *ReconciliationRepository) Create(rec *model.TransactionReconciliation) error {
	return r.db.Create(rec).Error
}

// GetByID 根据 ID 和用户 ID 获取单条对账记录。
// 通过 JOIN accounts 表验证该对账记录属于指定用户（权限校验）。
// 执行 SQL: SELECT * FROM transaction_reconciliations JOIN accounts ON ... WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 对账记录 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的对账记录对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *ReconciliationRepository) GetByID(id, userID uint64) (*model.TransactionReconciliation, error) {
	var rec model.TransactionReconciliation
	if err := r.db.Joins("JOIN accounts ON accounts.id = transaction_reconciliations.account_id").
		Where("transaction_reconciliations.id = ? AND accounts.user_id = ?", id, userID).
		First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

// List 分页获取指定用户的对账记录列表。
// 通过 JOIN accounts 表过滤出属于该用户的对账记录，可选按账户 ID 过滤。
// 执行 SQL: SELECT * FROM transaction_reconciliations JOIN accounts ON ... WHERE user_id = ? [AND account_id = ?] ORDER BY created_at DESC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 accountID: 可选的账户 ID 过滤条件，为 nil 时不按账户过滤。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 对账记录列表。
func (r *ReconciliationRepository) List(userID uint64, accountID *uint64, offset, limit int) ([]model.TransactionReconciliation, error) {
	var recs []model.TransactionReconciliation
	query := r.db.Model(&model.TransactionReconciliation{}).
		Joins("JOIN accounts ON accounts.id = transaction_reconciliations.account_id").
		Where("accounts.user_id = ?", userID)
	if accountID != nil {
		query = query.Where("transaction_reconciliations.account_id = ?", *accountID)
	}
	if err := query.Order("transaction_reconciliations.created_at DESC").Offset(offset).Limit(limit).Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

// Count 统计指定用户的对账记录总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM transaction_reconciliations JOIN accounts ON ... WHERE user_id = ? [AND account_id = ?]
func (r *ReconciliationRepository) Count(userID uint64, accountID *uint64) (int64, error) {
	var count int64
	query := r.db.Model(&model.TransactionReconciliation{}).
		Joins("JOIN accounts ON accounts.id = transaction_reconciliations.account_id").
		Where("accounts.user_id = ?", userID)
	if accountID != nil {
		query = query.Where("transaction_reconciliations.account_id = ?", *accountID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新对账记录。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE transaction_reconciliations SET ... WHERE id = ?
// 参数 rec: 要更新的对账记录对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *ReconciliationRepository) Update(rec *model.TransactionReconciliation) error {
	return r.db.Save(rec).Error
}

// Delete 删除对账记录及其关联的所有对账条目。
// 使用数据库事务确保原子性：先删除关联的对账条目，再删除对账记录本身。
// 如果任一步骤失败，整个事务回滚，数据保持一致。
// 执行 SQL（事务内）:
//   1. DELETE FROM reconciliation_entries WHERE reconciliation_id = ?
//   2. DELETE FROM transaction_reconciliations WHERE id = ?
// 参数 id: 要删除的对账记录 ID。
// 返回: 删除失败时返回错误。
func (r *ReconciliationRepository) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("reconciliation_id = ?", id).Delete(&model.ReconciliationEntry{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.TransactionReconciliation{}, id).Error
	})
}

// CreateEntry 创建一条对账条目。对账条目记录了对账中包含的具体交易信息。
// 执行 SQL: INSERT INTO reconciliation_entries (...)
// 参数 entry: 要创建的对账条目对象。
// 返回: 创建失败时返回错误。
func (r *ReconciliationRepository) CreateEntry(entry *model.ReconciliationEntry) error {
	return r.db.Create(entry).Error
}

// ListEntries 获取指定对账记录下的所有对账条目。
// 执行 SQL: SELECT * FROM reconciliation_entries WHERE reconciliation_id = ?
// 参数 reconciliationID: 对账记录 ID。
// 返回: 对账条目列表。
func (r *ReconciliationRepository) ListEntries(reconciliationID uint64) ([]model.ReconciliationEntry, error) {
	var entries []model.ReconciliationEntry
	if err := r.db.Where("reconciliation_id = ?", reconciliationID).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}
