// Package repository 提供数据访问层，封装所有与数据库的交互逻辑。
// 每个 Repository 结构体对应一个数据库表，提供增删改查等基本操作。
// 所有 Repository 都依赖 GORM 作为 ORM 框架，通过 *gorm.DB 执行数据库操作。
package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// AuthRepository 认证仓库，负责用户认证相关的数据访问。
// 与 UserRepository 不同，AuthRepository 专注于认证流程（注册、登录），
// 不需要 userID 参数进行数据隔离，因为认证时用户尚未登录或正在验证身份。
// 注意：用户管理相关的操作（列表、删除等）使用 UserRepository。
type AuthRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewAuthRepository 创建认证仓库实例。
// 参数 readDB: 读库（只读连接池，并发安全）
// 参数 writeDB: 写库（单连接，串行保证安全）
func NewAuthRepository(readDB, writeDB *gorm.DB) *AuthRepository {
	return &AuthRepository{readDB: readDB, writeDB: writeDB}
}

// Create 创建新用户记录，用于用户注册。
// 执行 SQL: INSERT INTO users (...)
// 参数 user: 要创建的用户对象，GORM 会自动填充 ID、CreatedAt 等字段。
// 返回: 创建失败时返回错误（如邮箱唯一约束冲突）。
func (r *AuthRepository) Create(user *model.User) error {
	return r.writeDB.Create(user).Error
}

// FindByEmail 根据邮箱查找用户，用于登录认证。
// 登录时用户输入邮箱和密码，系统先通过邮箱查找用户记录，再验证密码。
// 执行 SQL: SELECT * FROM users WHERE email = ? LIMIT 1
// 参数 email: 用户邮箱地址。
// 返回: 找到的用户对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *AuthRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.readDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息，用于修改密码、更新两步验证状态等认证相关操作。
// GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE users SET ... WHERE id = ?
// 参数 user: 要更新的用户对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *AuthRepository) Update(user *model.User) error {
	return r.writeDB.Save(user).Error
}

// Count 统计用户总数，用于判断是否为首个注册用户。
// 执行 SQL: SELECT COUNT(*) FROM users
// 返回: 用户总数和错误信息。
func (r *AuthRepository) Count() (int64, error) {
	var count int64
	if err := r.readDB.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetByID 根据 ID 获取用户，用于认证中间件从 token 中解析用户 ID 后获取用户信息。
// 注意：此方法不做用户权限校验，因为认证流程中用户身份尚未确认。
// 执行 SQL: SELECT * FROM users WHERE id = ? LIMIT 1
// 参数 id: 用户 ID。
// 返回: 找到的用户对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *AuthRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	if err := r.readDB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
