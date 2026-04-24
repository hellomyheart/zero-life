// Package service 业务逻辑层，实现核心业务逻辑
// AdminService 管理员业务逻辑，管理系统配置和用户管理
package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/email"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/hash"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// AdminService 管理员服务
// 负责系统管理功能，包括用户管理（列表、更新、删除、邀请）和系统配置管理
// 依赖userRepo进行用户数据访问，依赖configRepo进行系统配置数据访问
type AdminService struct {
	userRepo  *repository.UserRepository         // 用户数据访问对象
	configRepo *repository.ConfigurationRepository // 系统配置数据访问对象
}

// NewAdminService 创建管理员服务实例
// 参数：
//   - userRepo: 用户数据访问对象
//   - configRepo: 系统配置数据访问对象
// 返回：
//   - *AdminService: 管理员服务实例
func NewAdminService(
	userRepo *repository.UserRepository,
	configRepo *repository.ConfigurationRepository,
) *AdminService {
	return &AdminService{
		userRepo:  userRepo,
		configRepo: configRepo,
	}
}

// ListUsers 获取用户列表（分页）
// 支持关键词搜索，限制每页最多100条
// 参数：
//   - req: 列表查询参数（含搜索关键词、分页）
// 返回：
//   - []response.AdminUserResp: 用户列表
//   - int64: 总数
//   - error: 错误信息
func (s *AdminService) ListUsers(req *request.AdminListUsersReq) ([]response.AdminUserResp, int64, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	users, total, err := s.userRepo.List(req.Search, offset, pageSize)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}

	items := make([]response.AdminUserResp, 0, len(users))
	for _, u := range users {
		items = append(items, s.toUserResp(&u))
	}
	return items, total, nil
}

// UpdateUser 更新用户信息（管理员操作）
// 支持更新昵称、角色、语言、时区
// 参数：
//   - userID: 用户ID
//   - req: 更新请求参数
// 返回：
//   - *response.AdminUserResp: 更新后的用户信息
//   - error: 错误信息
func (s *AdminService) UpdateUser(userID uint64, req *request.AdminUpdateUserReq) (*response.AdminUserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.AdminUserResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Language:  user.Language,
		Timezone:  user.Timezone,
		MFAEnabled: user.MFAEnabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// DeleteUser 软删除用户
// 参数：
//   - userID: 用户ID
// 返回：
//   - error: 错误信息
func (s *AdminService) DeleteUser(userID uint64) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.userRepo.Delete(userID)
}

// InviteUser 邀请用户
// 创建新用户并生成随机临时密码，用户首次登录后需重置密码
// 参数：
//   - req: 邀请请求参数（邮箱、昵称、角色）
// 返回：
//   - *response.AdminUserResp: 创建的用户信息
//   - error: 错误信息（如邮箱已存在）
func (s *AdminService) InviteUser(req *request.AdminInviteUserReq) (*response.AdminUserResp, error) {
	// Check email uniqueness
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errcode.ErrEmailExists
	}

	// Generate a random temporary password
	tempPassword, err := hash.HashPassword("TempPass_" + randomHex(16))
	if err != nil {
		return nil, errcode.ErrInternal
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Email:    req.Email,
		Password: tempPassword,
		Nickname: req.Nickname,
		Role:     role,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.AdminUserResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Language:  user.Language,
		Timezone:  user.Timezone,
		MFAEnabled: user.MFAEnabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetConfiguration 获取系统配置项
// 参数：
//   - name: 配置名称
// 返回：
//   - *response.AdminConfigurationResp: 配置信息
//   - error: 错误信息
func (s *AdminService) GetConfiguration(name string) (*response.AdminConfigurationResp, error) {
	cfg, err := s.configRepo.Get(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return &response.AdminConfigurationResp{
		ID:        cfg.ID,
		Name:      cfg.Name,
		Value:     cfg.Value,
		CreatedAt: cfg.CreatedAt,
		UpdatedAt: cfg.UpdatedAt,
	}, nil
}

// ListConfigurations 获取所有系统配置
// 返回：
//   - []response.AdminConfigurationResp: 配置列表
//   - error: 错误信息
func (s *AdminService) ListConfigurations() ([]response.AdminConfigurationResp, error) {
	configs, err := s.configRepo.List()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.AdminConfigurationResp, 0, len(configs))
	for _, cfg := range configs {
		items = append(items, response.AdminConfigurationResp{
			ID:        cfg.ID,
			Name:      cfg.Name,
			Value:     cfg.Value,
			CreatedAt: cfg.CreatedAt,
			UpdatedAt: cfg.UpdatedAt,
		})
	}
	return items, nil
}

// UpdateConfiguration 更新系统配置项（不存在则创建）
// 参数：
//   - name: 配置名称
//   - req: 更新请求参数（含配置值）
// 返回：
//   - *response.AdminConfigurationResp: 更新后的配置信息
//   - error: 错误信息
func (s *AdminService) UpdateConfiguration(name string, req *request.AdminUpdateConfigurationReq) (*response.AdminConfigurationResp, error) {
	cfg := &model.Configuration{
		Name:  name,
		Value: req.Value,
	}
	if err := s.configRepo.Set(cfg); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.GetConfiguration(name)
}

// TestEmail 发送测试邮件，验证SMTP配置是否正确
// 参数：
//   - to: 收件人邮箱地址
// 返回：
//   - *response.AdminTestEmailResp: 测试结果（成功/失败及消息）
//   - error: 错误信息
func (s *AdminService) TestEmail(to string) (*response.AdminTestEmailResp, error) {
	err := email.Send(to, "Test Email from ZeroLife", `
	<html>
	<body>
		<h2>Test Email</h2>
		<p>This is a test email from ZeroLife to verify SMTP configuration.</p>
	</body>
	</html>`)
	if err != nil {
		return &response.AdminTestEmailResp{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &response.AdminTestEmailResp{
		Success: true,
		Message: "Test email sent successfully",
	}, nil
}

// toUserResp 将用户模型转换为管理员用户响应对象
func (s *AdminService) toUserResp(u *model.User) response.AdminUserResp {
	return response.AdminUserResp{
		ID:        u.ID,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Role:      u.Role,
		Language:  u.Language,
		Timezone:  u.Timezone,
		MFAEnabled: u.MFAEnabled,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// randomHex 生成指定长度的随机十六进制字符串
func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
