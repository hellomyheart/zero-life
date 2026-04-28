package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// RecurringTransactionRepository 定期交易仓库，负责定期交易及其执行日志的数据访问。
// 定期交易（RecurringTransaction）表示周期性自动创建的交易，如每月工资、每周通勤费等。
// 每次执行时会创建一条日志记录（RecurringTransactionLog），并更新下次执行时间。
type RecurringTransactionRepository struct {
	db *gorm.DB
}

// NewRecurringTransactionRepository 创建定期交易仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewRecurringTransactionRepository(db *gorm.DB) *RecurringTransactionRepository {
	return &RecurringTransactionRepository{db: db}
}

// Create 创建一条新的定期交易记录。
// 执行 SQL: INSERT INTO recurring_transactions (...)
// 参数 rt: 要创建的定期交易对象。
// 返回: 创建失败时返回错误。
func (r *RecurringTransactionRepository) Create(rt *model.RecurringTransaction) error {
	return r.db.Create(rt).Error
}

// GetByID 根据 ID 和用户 ID 获取单条定期交易。
// 执行 SQL: SELECT * FROM recurring_transactions WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 定期交易 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的定期交易对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *RecurringTransactionRepository) GetByID(id, userID uint64) (*model.RecurringTransaction, error) {
	var rt model.RecurringTransaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// List 分页获取指定用户的定期交易列表，可选按激活状态过滤。
// 执行 SQL: SELECT * FROM recurring_transactions WHERE user_id = ? [AND is_active = ?] ORDER BY next_occurrence ASC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 active: 可选的激活状态过滤，为 nil 时不按状态过滤。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 定期交易列表。
func (r *RecurringTransactionRepository) List(userID uint64, active *bool, offset, limit int) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	query := r.db.Where("user_id = ?", userID)
	if active != nil {
		query = query.Where("is_active = ?", *active)
	}
	if err := query.Order("next_occurrence ASC").Offset(offset).Limit(limit).Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

// Count 统计指定用户的定期交易总数，可选按激活状态过滤，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM recurring_transactions WHERE user_id = ? [AND is_active = ?]
func (r *RecurringTransactionRepository) Count(userID uint64, active *bool) (int64, error) {
	var count int64
	query := r.db.Model(&model.RecurringTransaction{}).Where("user_id = ?", userID)
	if active != nil {
		query = query.Where("is_active = ?", *active)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新定期交易。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE recurring_transactions SET ... WHERE id = ?
// 参数 rt: 要更新的定期交易对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *RecurringTransactionRepository) Update(rt *model.RecurringTransaction) error {
	return r.db.Save(rt).Error
}

func (r *RecurringTransactionRepository) UpdateWithDB(dbTx *gorm.DB, rt *model.RecurringTransaction) error {
	return dbTx.Save(rt).Error
}

// Delete 删除定期交易及其关联的所有执行日志。使用数据库事务确保原子性。
// 执行 SQL（事务内）:
//   1. DELETE FROM recurring_transaction_logs WHERE recurring_transaction_id = ?
//   2. DELETE FROM recurring_transactions WHERE id = ? AND user_id = ?
// 参数 id: 定期交易 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误，事务回滚。
func (r *RecurringTransactionRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("recurring_transaction_id = ?", id).Delete(&model.RecurringTransactionLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecurringTransaction{}).Error
	})
}

// GetDueRecurringTransactions 获取指定用户中已到期的定期交易，用于定时任务处理。
// 使用数据库的 NOW() 函数比较，确保时间判断在数据库端完成。
// 执行 SQL: SELECT * FROM recurring_transactions WHERE user_id = ? AND is_active = true AND next_occurrence <= NOW()
// 参数 userID: 用户 ID。
// 返回: 已到期的定期交易列表。
func (r *RecurringTransactionRepository) GetDueRecurringTransactions(userID uint64) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := gorm.Expr("NOW()")
	if err := r.db.Where("user_id = ? AND is_active = ? AND next_occurrence <= ?", userID, true, now).
		Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

// CreateLog 创建一条定期交易执行日志，记录每次自动创建交易的情况。
// 执行 SQL: INSERT INTO recurring_transaction_logs (...)
// 参数 log: 要创建的日志对象。
// 返回: 创建失败时返回错误。
func (r *RecurringTransactionRepository) CreateLog(log *model.RecurringTransactionLog) error {
	return r.db.Create(log).Error
}

func (r *RecurringTransactionRepository) CreateLogWithDB(dbTx *gorm.DB, log *model.RecurringTransactionLog) error {
	return dbTx.Create(log).Error
}

// ListLogs 获取指定定期交易的所有执行日志，按执行日期倒序排列（最新的在前）。
// 执行 SQL: SELECT * FROM recurring_transaction_logs WHERE recurring_transaction_id = ? ORDER BY occurrence_date DESC
// 参数 recurringTransactionID: 定期交易 ID。
// 返回: 执行日志列表。
func (r *RecurringTransactionRepository) ListLogs(recurringTransactionID uint64) ([]model.RecurringTransactionLog, error) {
	var logs []model.RecurringTransactionLog
	if err := r.db.Where("recurring_transaction_id = ?", recurringTransactionID).
		Order("occurrence_date DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetUpcoming 获取指定用户未来若干天内到期的循环交易，用于仪表盘提醒。
// 仅返回 reminder_days > 0 且 is_active = true 的记录。
// 提前提醒逻辑：到期前 reminder_days 天就开始显示，即 next_occurrence - reminder_days 天 <= now + days。
func (r *RecurringTransactionRepository) GetUpcoming(userID uint64, days int) ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := time.Now()
	endDate := now.AddDate(0, 0, days)
	// SQLite: date(next_occurrence, '-' || reminder_days || ' days') 计算提醒开始日期
	if err := r.db.Where(
		"user_id = ? AND is_active = ? AND reminder_days > 0 AND next_occurrence <= ? AND date(next_occurrence, '-' || reminder_days || ' days') <= ?",
		userID, true, endDate, endDate.Format("2006-01-02")).
		Order("next_occurrence ASC").Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

// GetAllDue 获取所有用户的到期循环交易（不按用户过滤），用于定时任务批量处理。
func (r *RecurringTransactionRepository) GetAllDue() ([]model.RecurringTransaction, error) {
	var rts []model.RecurringTransaction
	now := time.Now()
	if err := r.db.Where("is_active = ? AND next_occurrence <= ?", true, now).
		Order("next_occurrence ASC").Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}
