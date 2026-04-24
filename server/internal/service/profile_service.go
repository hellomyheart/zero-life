// Package service 业务逻辑层，实现核心业务逻辑
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ProfileService 用户资料业务服务
// 负责处理用户资料相关的业务逻辑
type ProfileService struct {
	userRepo *repository.UserRepository
}

// NewProfileService 创建用户资料服务实例
// 参数：
//   - userRepo: 用户数据访问对象
// 返回：
//   - *ProfileService: 用户资料服务实例
func NewProfileService(userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

// GetProfile 获取用户资料
// 参数：
//   - userID: 用户ID
// 返回：
//   - *response.UserResp: 用户资料信息
//   - error: 错误信息
func (s *ProfileService) GetProfile(userID uint64) (*response.UserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(user), nil
}

// UpdateProfile 更新用户资料
// 参数：
//   - userID: 用户ID
//   - req: 更新请求参数
// 返回：
//   - *response.UserResp: 更新后的用户资料
//   - error: 错误信息
func (s *ProfileService) UpdateProfile(userID uint64, req *request.UpdateProfileReq) (*response.UserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(user), nil
}

// ChangePassword 修改密码
// 参数：
//   - userID: 用户ID
//   - req: 修改密码请求参数
// 返回：
//   - error: 错误信息
func (s *ProfileService) ChangePassword(userID uint64, req *request.ChangePasswordReq) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errcode.ErrInvalidPassword
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternal
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// toResp 将用户模型转换为响应对象
// 参数：
//   - user: 用户模型
// 返回：
//   - *response.UserResp: 用户响应对象
func (s *ProfileService) toResp(user *model.User) *response.UserResp {
	return &response.UserResp{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
