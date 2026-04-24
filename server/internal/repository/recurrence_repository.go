package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// RecurrenceRepository 周期性事务仓库，负责周期性事务的数据访问。
// 周期性事务（Recurrence）与 RecurringTransaction 类似，表示周期性重复的财务事件。
// 包含下次执行日期（next_date）和激活状态（is_active），定时任务会检查到期的事务。
type RecurrenceRepository struct {
	db *gorm.DB
}

// NewRecurrenceRepository 创建周期性事务仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewRecurrenceRepository(db *gorm.DB) *RecurrenceRepository {
	return &RecurrenceRepository{db: db}
}

// Create 创建一条新的周期性事务记录。
// 执行 SQL: INSERT INTO recurrences (...)
// 参数 rec: 要创建的周期性事务对象。
// 返回: 创建失败时返回错误。
func (r *RecurrenceRepository) Create(rec *model.Recurrence) error {
	return r.db.Create(rec).Error
}

// GetByID 根据 ID 和用户 ID 获取单条周期性事务。
// 执行 SQL: SELECT * FROM recurrences WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 周期性事务 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的周期性事务对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *RecurrenceRepository) GetByID(id, userID uint64) (*model.Recurrence, error) {
	var rec model.Recurrence
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

// List 分页获取指定用户的周期性事务列表，按下次执行日期升序排列。
// 执行 SQL: SELECT * FROM recurrences WHERE user_id = ? ORDER BY next_date ASC LIMIT ? OFFSET ?
// 参数 userID: 用户 ID。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 周期性事务列表。
func (r *RecurrenceRepository) List(userID uint64, offset, limit int) ([]model.Recurrence, error) {
	var recs []model.Recurrence
	if err := r.db.Where("user_id = ?", userID).Order("next_date ASC").Offset(offset).Limit(limit).Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

// Count 统计指定用户的周期性事务总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM recurrences WHERE user_id = ?
func (r *RecurrenceRepository) Count(userID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Recurrence{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新周期性事务。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE recurrences SET ... WHERE id = ?
// 参数 rec: 要更新的周期性事务对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *RecurrenceRepository) Update(rec *model.Recurrence) error {
	return r.db.Save(rec).Error
}

// Delete 根据 ID 和用户 ID 删除周期性事务。
// 执行 SQL: DELETE FROM recurrences WHERE id = ? AND user_id = ?
// 参数 id: 周期性事务 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误。
func (r *RecurrenceRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Recurrence{}).Error
}

// GetDueRecurrences 获取所有已到期的激活周期性事务，用于定时任务处理。
// 不按用户过滤，因为定时任务需要处理所有用户的到期事务。
// 执行 SQL: SELECT * FROM recurrences WHERE is_active = true AND next_date <= ? ORDER BY next_date ASC
// 参数 today: 当前日期，用于判断是否到期。
// 返回: 已到期的周期性事务列表。
func (r *RecurrenceRepository) GetDueRecurrences(today time.Time) ([]model.Recurrence, error) {
	var recs []model.Recurrence
	if err := r.db.Where("is_active = ? AND next_date <= ?", true, today).
		Order("next_date ASC").Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}
