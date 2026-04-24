package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓库，负责用户管理的数据访问。
// 提供用户的增删改查、搜索、邮箱查找等功能，主要用于管理员管理用户。
// 注意：认证相关的用户操作（注册、登录）使用 AuthRepository。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// List 分页获取用户列表，支持按邮箱或昵称搜索。
// 先统计总数，再查询分页数据。搜索时使用 LIKE 模糊匹配。
// 执行 SQL:
//   1. SELECT COUNT(*) FROM users WHERE email LIKE ? OR nickname LIKE ?
//   2. SELECT * FROM users WHERE email LIKE ? OR nickname LIKE ? ORDER BY id ASC LIMIT ? OFFSET ?
// 参数 search: 搜索关键词，为空时不进行过滤。同时匹配邮箱和昵称。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 用户列表、总记录数、错误信息。
func (r *UserRepository) List(search string, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("email LIKE ? OR nickname LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// GetByID 根据 ID 获取单个用户。
// 执行 SQL: SELECT * FROM users WHERE id = ? LIMIT 1
// 参数 id: 用户 ID。
// 返回: 找到的用户对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE users SET ... WHERE id = ?
// 参数 user: 要更新的用户对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 根据 ID 删除用户。
// 注意：此方法不做软删除，而是直接从数据库中移除记录。
// 执行 SQL: DELETE FROM users WHERE id = ?
// 参数 id: 要删除的用户 ID。
// 返回: 删除失败时返回错误。
func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&model.User{}, id).Error
}

// FindByEmail 根据邮箱查找用户，用于登录认证等场景。
// 执行 SQL: SELECT * FROM users WHERE email = ? LIMIT 1
// 参数 email: 用户邮箱地址。
// 返回: 找到的用户对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建新用户。
// 执行 SQL: INSERT INTO users (...)
// 参数 user: 要创建的用户对象。
// 返回: 创建失败时返回错误。
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Count 获取用户总数，用于统计和管理面板展示。
// 执行 SQL: SELECT COUNT(*) FROM users
// 返回: 用户总数。
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}
