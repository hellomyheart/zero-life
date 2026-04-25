// Package service 业务逻辑层，实现核心业务逻辑
// UserService 用户业务逻辑，处理用户管理和角色权限
package service

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
// 处理用户管理相关业务逻辑，包括用户列表、详情、更新、删除、角色管理、锁定/解锁和密码重置
// 依赖userRepo查询用户数据，依赖db执行直接更新操作（角色、锁定、密码等）
type UserService struct {
	userRepo *repository.UserRepository // 用户数据访问对象
	db       *gorm.DB                   // 数据库连接，用于直接更新用户字段
}

// NewUserService 创建用户服务实例
// 参数：
//   - userRepo: 用户数据访问对象
//   - db: 数据库连接
// 返回：
//   - *UserService: 用户服务实例
func NewUserService(userRepo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
	}
}

// List 获取用户列表（分页）
// 参数：
//   - req: 列表查询参数（含分页）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *UserService) List(req *request.UserListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	users, total, err := s.userRepo.List("", params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.UserResp, 0, len(users))
	for _, u := range users {
		items = append(items, *s.toResp(&u))
	}

	return pagination.NewResult(items, total, params), nil
}

// Get 获取用户详情
// 参数：
//   - id: 用户ID
// 返回：
//   - *response.UserResp: 用户信息
//   - error: 错误信息（如用户不存在）
func (s *UserService) Get(id uint64) (*response.UserResp, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	return s.toResp(user), nil
}

// Update 更新用户信息
// 支持更新昵称、语言、时区（部分更新）
// 参数：
//   - id: 用户ID
//   - req: 更新请求参数
// 返回：
//   - *response.UserResp: 更新后的用户信息
//   - error: 错误信息
func (s *UserService) Update(id uint64, req *request.UpdateUserReq) (*response.UserResp, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(user), nil
}

// Delete 删除用户（软删除）
// 参数：
//   - id: 用户ID
// 返回：
//   - error: 错误信息
func (s *UserService) Delete(id uint64) error {
	// 软删除
	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		return errcode.ErrInternal
	}
	return nil
}

// ChangeRole 修改用户角色
// 业务规则：角色只能是"user"或"admin"，其他值返回错误
// 参数：
//   - id: 用户ID
//   - role: 新角色（user/admin）
// 返回：
//   - error: 错误信息（如角色无效）
func (s *UserService) ChangeRole(id uint64, role string) error {
	// 验证角色
	if role != "user" && role != "admin" {
		return errcode.WithMessage(errcode.ErrBadRequest, "invalid role")
	}

	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"role":       role,
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Lock 锁定用户
// 将用户标记为锁定状态，锁定后用户无法登录
// 锁定不改变用户角色，解锁后保留原有角色权限
// 参数：
//   - id: 用户ID
// 返回：
//   - error: 错误信息
func (s *UserService) Lock(id uint64) error {
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_locked":  true,
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Unlock 解锁用户
// 将用户标记为未锁定状态，恢复登录权限
// 解锁后用户保留原有角色，不会降级为普通用户
// 参数：
//   - id: 用户ID
// 返回：
//   - error: 错误信息
func (s *UserService) Unlock(id uint64) error {
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_locked":  false,
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// ResetPassword 重置用户密码（管理员操作）
// 业务流程：
// 1. 验证新密码长度（至少8位）
// 2. 使用bcrypt加密密码
// 3. 更新用户密码
// 参数：
//   - id: 用户ID
//   - newPassword: 新密码
// 返回：
//   - error: 错误信息（如密码太短）
func (s *UserService) ResetPassword(id uint64, newPassword string) error {
	// 验证密码长度
	if len(newPassword) < 8 {
		return errcode.ErrPasswordTooShort
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternal
	}

	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password":   string(hashedPassword),
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// toResp 将用户模型转换为响应对象
// 参数：
//   - u: 用户模型
// 返回：
//   - *response.UserResp: 用户响应对象
func (s *UserService) toResp(u *model.User) *response.UserResp {
	return &response.UserResp{
		ID:         u.ID,
		Email:      u.Email,
		Nickname:   u.Nickname,
		Language:   u.Language,
		Timezone:   u.Timezone,
		Role:       u.Role,
		MFAEnabled: u.MFAEnabled,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
