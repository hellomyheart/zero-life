// Package service 业务逻辑层
// 实现核心业务逻辑，调用Repository层进行数据操作
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
// 处理用户管理相关业务逻辑
type UserService struct {
	userRepo *repository.UserRepository
	db       *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
	}
}

// List 获取用户列表
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

// Delete 删除用户
func (s *UserService) Delete(id uint64) error {
	// 软删除
	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		return errcode.ErrInternal
	}
	return nil
}

// ChangeRole 修改用户角色
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
func (s *UserService) Lock(id uint64) error {
	// 这里可以通过设置一个锁定标记或修改角色来实现
	// 简单实现：将角色改为 "locked"
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"role":       "locked",
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Unlock 解锁用户
func (s *UserService) Unlock(id uint64) error {
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"role":       "user",
			"updated_at": now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// ResetPassword 重置用户密码（管理员操作）
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

// toResp 转换为响应格式
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
