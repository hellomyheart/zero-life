package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// PreferenceRepository 偏好设置仓库，负责用户偏好设置的数据访问。
// 偏好设置（Preference）以键值对形式存储用户的个性化配置，
// 如默认货币、界面语言、日期格式等。每个用户可以有多组键值对。
type PreferenceRepository struct {
	db *gorm.DB
}

// NewPreferenceRepository 创建偏好设置仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

// Get 获取指定用户的某个偏好设置值。
// 执行 SQL: SELECT * FROM preferences WHERE user_id = ? AND key = ? LIMIT 1
// 参数 userID: 用户 ID。
// 参数 key: 偏好设置的键名，如 "currency"、"language" 等。
// 返回: 找到的偏好设置对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *PreferenceRepository) Get(userID uint64, key string) (*model.Preference, error) {
	var pref model.Preference
	if err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

// Set 设置指定用户的偏好值，实现"存在则更新，不存在则创建"的逻辑（Upsert）。
// 先查询该键是否存在，如果不存在（ErrRecordNotFound）则创建新记录，否则更新已有记录。
// 执行 SQL:
//   查询: SELECT * FROM preferences WHERE user_id = ? AND key = ? LIMIT 1
//   创建: INSERT INTO preferences (user_id, key, value) VALUES (?, ?, ?)
//   更新: UPDATE preferences SET value = ? WHERE id = ?
// 参数 userID: 用户 ID。
// 参数 key: 偏好设置的键名。
// 参数 value: 偏好设置的值。
// 返回: 操作失败时返回错误。
func (r *PreferenceRepository) Set(userID uint64, key, value string) error {
	var pref model.Preference
	err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		pref = model.Preference{UserID: userID, Key: key, Value: value}
		return r.db.Create(&pref).Error
	}
	if err != nil {
		return err
	}
	pref.Value = value
	return r.db.Save(&pref).Error
}

// List 获取指定用户的所有偏好设置。
// 执行 SQL: SELECT * FROM preferences WHERE user_id = ?
// 参数 userID: 用户 ID。
// 返回: 偏好设置列表。
func (r *PreferenceRepository) List(userID uint64) ([]model.Preference, error) {
	var prefs []model.Preference
	if err := r.db.Where("user_id = ?", userID).Find(&prefs).Error; err != nil {
		return nil, err
	}
	return prefs, nil
}

// Delete 删除指定用户的某个偏好设置。
// 执行 SQL: DELETE FROM preferences WHERE user_id = ? AND key = ?
// 参数 userID: 用户 ID。
// 参数 key: 要删除的偏好设置键名。
// 返回: 删除失败时返回错误。
func (r *PreferenceRepository) Delete(userID uint64, key string) error {
	return r.db.Where("user_id = ? AND key = ?", userID, key).Delete(&model.Preference{}).Error
}
