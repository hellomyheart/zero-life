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
	"github.com/hellomyheart/zero-life/server/internal/pkg/totp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	loginFailKeyPrefix  = "login_fail:"    // 登录失败次数键前缀
	loginLockKeyPrefix  = "login_lock:"    // 登录锁定键前缀
	resetTokenKeyPrefix = "reset_token:"   // 密码重置令牌键前缀
	mfaTokenKeyPrefix   = "mfa_token:"     // MFA登录验证令牌键前缀
	maxLoginFails       = 5                // 最大登录失败次数
	loginLockDuration   = 30 * time.Minute // 登录锁定时长
	resetTokenTTL       = 24 * time.Hour   // 密码重置令牌有效期
	mfaTokenTTL         = 5 * time.Minute  // MFA登录验证令牌有效期
)

// AuthService 认证服务
// 处理用户注册、登录、密码重置等认证相关业务逻辑
type AuthService struct {
	authRepo   *repository.AuthRepository
	jwtService *jwt.Service
	kvRepo     *repository.KVRepository
	mfaService *MFAService
}

func NewAuthService(authRepo *repository.AuthRepository, jwtService *jwt.Service, kvRepo *repository.KVRepository, mfaService *MFAService) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtService: jwtService,
		kvRepo:     kvRepo,
		mfaService: mfaService,
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

	count, err := s.authRepo.Count()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	if count == 0 {
		user.Role = "admin"
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
// 1. 根据邮箱查找用户
// 2. 检查账户是否被管理员锁定（user.IsLocked）
// 3. 检查账户是否被临时锁定（kv_store 中存在锁定键则拒绝登录）
// 4. 验证密码是否正确
// 5. 密码错误时：记录失败次数到 kv_store，连续失败5次则锁定账户30分钟
// 6. 密码正确时：清除失败计数，生成JWT令牌对
// 参数：
//   - req: 登录请求（邮箱、密码）
// 返回：
//   - *response.LoginResp: 登录响应（包含访问令牌和刷新令牌）
//   - error: 错误信息（如账户锁定、凭证错误）
func (s *AuthService) Login(req *request.LoginReq) (*response.LoginResp, error) {
	// 步骤1：根据邮箱查找用户
	user, err := s.authRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrInvalidCredential
		}
		return nil, errcode.ErrInternal
	}

	// 步骤2：检查管理员手动锁定（持久锁定）
	if user.IsLocked {
		return nil, errcode.ErrAccountLocked
	}

	// 步骤3：检查临时锁定（kv_store 连续失败5次锁定）
	lockKey := loginLockKeyPrefix + req.Email
	locked, err := s.kvRepo.Exists(lockKey)
	if err == nil && locked {
		return nil, errcode.ErrAccountLocked
	}

	// 步骤4：验证密码
	if !hash.CheckPassword(req.Password, user.Password) {
	// 步骤4a：密码错误，记录失败次数
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

	// 步骤5：密码正确，清除失败计数
	failKey := loginFailKeyPrefix + req.Email
	s.kvRepo.Del(failKey)

	// 步骤6：如果用户启用了MFA，生成临时验证令牌，不返回真正的TokenPair
	if user.MFAEnabled {
		mfaTokenBytes := make([]byte, 32)
		if _, err := rand.Read(mfaTokenBytes); err != nil {
			return nil, errcode.ErrInternal
		}
		mfaToken := hex.EncodeToString(mfaTokenBytes)
		mfaKey := mfaTokenKeyPrefix + mfaToken
		if err := s.kvRepo.Set(mfaKey, fmt.Sprintf("%d", user.ID), mfaTokenTTL); err != nil {
			return nil, errcode.ErrInternal
		}
		return &response.LoginResp{
			MFARequired: true,
			AccessToken:  mfaToken,
			ExpiresAt:    time.Now().Add(mfaTokenTTL),
		}, nil
	}

	// 步骤7：未启用MFA，直接生成JWT令牌对
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

// MFALoginVerify MFA登录二次验证
// 用户启用MFA后，登录时密码正确会返回临时mfa_token，前端引导用户输入TOTP码后调用此方法
// 业务流程：
// 1. 从kv_store中验证mfa_token，获取用户ID
// 2. 查找用户并验证MFA状态
// 3. 验证TOTP码是否正确
// 4. 验证通过后删除临时令牌，生成真正的JWT令牌对
// 参数：
//   - req: MFA验证请求（含临时令牌和TOTP码）
// 返回：
//   - *response.LoginResp: 登录响应（含真正的访问令牌和刷新令牌）
//   - error: 错误信息
func (s *AuthService) MFALoginVerify(req *request.MFALoginVerifyReq) (*response.LoginResp, error) {
	mfaKey := mfaTokenKeyPrefix + req.MFAToken

	userIDStr, err := s.kvRepo.Get(mfaKey)
	if err != nil {
		return nil, errcode.ErrMFAInvalidToken
	}

	var userID uint64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		s.kvRepo.Del(mfaKey)
		return nil, errcode.ErrMFAInvalidToken
	}

	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		s.kvRepo.Del(mfaKey)
		return nil, errcode.ErrMFAInvalidToken
	}

	if !user.MFAEnabled || user.MFASecret == "" {
		s.kvRepo.Del(mfaKey)
		return nil, errcode.ErrMFANotEnabled
	}

	if !totp.ValidateCode(user.MFASecret, req.Code) {
		// TOTP验证失败，尝试备用码
		if !s.mfaService.VerifyBackupCode(userID, req.Code) {
			return nil, errcode.ErrMFAInvalidCode
		}
	}

	s.kvRepo.Del(mfaKey)

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.LoginResp{
		MFARequired:  false,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}