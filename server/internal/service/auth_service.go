// Package service 业务逻辑层
// 实现核心业务逻辑，包括认证、账户、交易等业务处理
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/config"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/email"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/hash"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	loginFailKeyPrefix  = "login_fail:"    // 登录失败次数键前缀
	loginLockKeyPrefix  = "login_lock:"    // 登录锁定键前缀
	resetTokenKeyPrefix = "reset_token:"   // 密码重置令牌键前缀
	maxLoginFails       = 5                // 最大登录失败次数
	loginLockDuration   = 30 * time.Minute // 登录锁定时长
	resetTokenTTL       = 24 * time.Hour   // 密码重置令牌有效期
)

// AuthService 认证服务
// 处理用户注册、登录、密码重置等认证相关业务逻辑
type AuthService struct {
	authRepo   *repository.AuthRepository // 用户仓储
	jwtService *jwt.Service               // JWT服务
	kvRepo     *repository.KVRepository   // 键值存储（用于登录限制和令牌管理，替代 Redis）
}

// NewAuthService 创建认证服务实例
// 参数：
//   authRepo: 用户仓储
//   jwtService: JWT服务
//   kvRepo: 键值存储仓库（替代 Redis）
func NewAuthService(authRepo *repository.AuthRepository, jwtService *jwt.Service, kvRepo *repository.KVRepository) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtService: jwtService,
		kvRepo:     kvRepo,
	}
}

// Register 用户注册
// 流程：
// 1. 检查邮箱是否已存在
// 2. 加密密码
// 3. 创建用户记录
// 4. 生成JWT令牌对
// 参数：
//   req: 注册请求（邮箱、密码、昵称）
// 返回：
//   LoginResp: 登录响应（包含访问令牌和刷新令牌）
//   error: 错误信息
func (s *AuthService) Register(req *request.RegisterReq) (*response.LoginResp, error) {
	existing, err := s.authRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errcode.ErrEmailExists
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errcode.ErrInternal
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	user := &model.User{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
	}
	if err := s.authRepo.Create(user); err != nil {
		return nil, errcode.ErrInternal
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// Login 用户登录
// 业务流程：
// 1. 检查账户是否被锁定（kv_store 中存在锁定键则拒绝登录）
// 2. 根据邮箱查找用户
// 3. 验证密码是否正确
// 4. 密码错误时：记录失败次数到 kv_store，连续失败5次则锁定账户30分钟
// 5. 密码正确时：清除失败计数，生成JWT令牌对
// 参数：
//   - req: 登录请求（邮箱、密码）
// 返回：
//   - *response.LoginResp: 登录响应（包含访问令牌和刷新令牌）
//   - error: 错误信息（如账户锁定、凭证错误）
func (s *AuthService) Login(req *request.LoginReq) (*response.LoginResp, error) {
	// 步骤1：检查账户是否被锁定
	lockKey := loginLockKeyPrefix + req.Email
	locked, err := s.kvRepo.Exists(lockKey)
	if err == nil && locked {
		return nil, errcode.ErrAccountLocked
	}

	// 步骤2：根据邮箱查找用户
	user, err := s.authRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrInvalidCredential
		}
		return nil, errcode.ErrInternal
	}

	// 步骤3：验证密码
	if !hash.CheckPassword(req.Password, user.Password) {
		// 步骤3a：密码错误，记录失败次数
		failKey := loginFailKeyPrefix + req.Email
		count, _ := s.kvRepo.Incr(failKey)
		if count == 1 {
			// 首次失败，设置失败计数器的过期时间（与锁定时长相同）
			s.kvRepo.Expire(failKey, loginLockDuration)
		}
		if count >= maxLoginFails {
			// 连续失败5次，锁定账户30分钟，并清除失败计数器
			s.kvRepo.Set(lockKey, "1", loginLockDuration)
			s.kvRepo.Del(failKey)
			return nil, errcode.ErrAccountLocked
		}
		return nil, errcode.ErrInvalidCredential
	}

	// 步骤4：密码正确，清除失败计数
	failKey := loginFailKeyPrefix + req.Email
	s.kvRepo.Del(failKey)

	// 步骤5：生成JWT令牌对
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// RefreshToken 刷新访问令牌
// 使用刷新令牌获取新的访问令牌对
// 业务流程：
// 1. 解析刷新令牌获取用户信息
// 2. 验证用户仍然存在（防止已删除用户的令牌继续使用）
// 3. 生成新的JWT令牌对
// 参数：
//   - refreshToken: 刷新令牌
// 返回：
//   - *response.LoginResp: 新的令牌对
//   - error: 错误信息（如令牌无效）
func (s *AuthService) RefreshToken(refreshToken string) (*response.LoginResp, error) {
	claims, err := s.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}

	user, err := s.authRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// ForgotPassword 忘记密码
// 业务流程：
// 1. 查找用户（用户不存在时不暴露此信息，直接返回nil防止邮箱枚举攻击）
// 2. 生成32字节随机令牌
// 3. 将令牌存入 kv_store，有效期24小时，值为用户ID
// 4. 发送包含重置链接的邮件（发送失败不阻断请求，令牌仍有效）
// 参数：
//   - emailAddr: 用户邮箱地址
// 返回：
//   - error: 错误信息（注意：用户不存在时也返回nil）
func (s *AuthService) ForgotPassword(emailAddr string) error {
	user, err := s.authRepo.FindByEmail(emailAddr)
	if err != nil {
		// Don't reveal whether email exists
		return nil
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return errcode.ErrInternal
	}
	token := hex.EncodeToString(tokenBytes)

	// Store in kv_store with 24h TTL
	key := resetTokenKeyPrefix + token
	if err := s.kvRepo.Set(key, fmt.Sprintf("%d", user.ID), resetTokenTTL); err != nil {
		return errcode.ErrInternal
	}

	// Send email with reset link
	frontendURL := config.C.App.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	if err := email.SendPasswordReset(user.Email, resetURL); err != nil {
		zap.L().Error("failed to send password reset email", zap.Error(err), zap.String("email", user.Email))
		// Don't fail the request - token is still stored
	}

	return nil
}

// ResetPassword 重置密码
// 业务流程：
// 1. 从 kv_store 中验证重置令牌，获取用户ID
// 2. 验证用户存在
// 3. 加密新密码并更新用户记录
// 4. 删除已使用的令牌（防止重复使用）
// 参数：
//   - token: 密码重置令牌（由ForgotPassword生成）
//   - newPassword: 新密码
// 返回：
//   - error: 错误信息（如令牌无效、用户不存在）
func (s *AuthService) ResetPassword(token, newPassword string) error {
	key := resetTokenKeyPrefix + token

	// Verify token
	userIDStr, err := s.kvRepo.Get(key)
	if err != nil {
		return errcode.ErrInvalidToken
	}

	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)

	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errcode.ErrInternal
	}

	user.Password = hashedPassword
	if err := s.authRepo.Update(user); err != nil {
		return errcode.ErrInternal
	}

	// Delete used token
	s.kvRepo.Del(key)

	return nil
}

// GetProfile 获取用户个人资料
// 参数：
//   - userID: 用户ID
// 返回：
//   - *response.ProfileResp: 用户资料信息
//   - error: 错误信息
func (s *AuthService) GetProfile(userID uint64) (*response.ProfileResp, error) {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	return &response.ProfileResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Language:  user.Language,
		Timezone:  user.Timezone,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateProfile 更新用户个人资料
// 支持更新昵称、语言、时区（部分更新，仅更新提供的字段）
// 参数：
//   - userID: 用户ID
//   - req: 更新请求参数
// 返回：
//   - *response.ProfileResp: 更新后的用户资料
//   - error: 错误信息
func (s *AuthService) UpdateProfile(userID uint64, req *request.UpdateProfileReq) (*response.ProfileResp, error) {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}

	if err := s.authRepo.Update(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.ProfileResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Language:  user.Language,
		Timezone:  user.Timezone,
		CreatedAt: user.CreatedAt,
	}, nil
}

// ChangePassword 修改密码
// 业务流程：
// 1. 验证旧密码是否正确
// 2. 加密新密码
// 3. 更新用户密码
// 参数：
//   - userID: 用户ID
//   - req: 修改密码请求（含旧密码和新密码）
// 返回：
//   - error: 错误信息（如旧密码错误）
func (s *AuthService) ChangePassword(userID uint64, req *request.ChangePasswordReq) error {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrNotFound
	}

	if !hash.CheckPassword(req.OldPassword, user.Password) {
		return errcode.ErrOldPasswordWrong
	}

	hashedPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return errcode.ErrInternal
	}

	user.Password = hashedPassword
	return s.authRepo.Update(user)
}