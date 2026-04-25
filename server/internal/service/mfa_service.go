// Package service 业务逻辑层
// 实现核心业务逻辑，调用Repository层进行数据操作
package service

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

// MFAService MFA（多因素认证）服务
// 处理多因素认证相关业务逻辑，基于TOTP（基于时间的一次性密码）算法
// 依赖userRepo查询用户信息，依赖db直接更新用户MFA字段
type MFAService struct {
	userRepo *repository.UserRepository // 用户数据访问对象
	db       *gorm.DB                   // 数据库连接，用于直接更新MFA字段
}

// NewMFAService 创建MFA服务实例
// 参数：
//   - userRepo: 用户数据访问对象
//   - db: 数据库连接
// 返回：
//   - *MFAService: MFA服务实例
func NewMFAService(userRepo *repository.UserRepository, db *gorm.DB) *MFAService {
	return &MFAService{
		userRepo: userRepo,
		db:       db,
	}
}

// Setup 初始化MFA设置
// 业务流程：
// 1. 获取用户信息
// 2. 检查是否已启用MFA（已启用则不能再初始化）
// 3. 生成随机MFA密钥（20字节随机数的Base32编码）
// 4. 生成TOTP二维码URL（供认证器App扫描）
// 5. 临时保存密钥到数据库（此时MFA未启用，需通过Enable方法验证后启用）
// 参数：
//   - userID: 用户ID
// 返回：
//   - *response.MFASetupResp: MFA设置信息（含密钥和二维码URL）
//   - error: 错误信息
func (s *MFAService) Setup(userID uint64) (*response.MFASetupResp, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// 检查是否已启用MFA
	if user.MFAEnabled {
		return nil, errcode.ErrMFAAlreadyEnabled
	}

	// 生成随机密钥
	secret, err := generateMFASecret()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// 生成二维码URL
	issuer := "ZeroLife"
	accountName := user.Email
	otpURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, accountName, secret, issuer)

	// 临时保存密钥到数据库（未启用状态）
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"mfa_secret": secret,
			"updated_at": now,
		}).Error; err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.MFASetupResp{
		Secret: secret,
		QRCode: otpURL,
	}, nil
}

// Enable 启用MFA
// 业务流程：
// 1. 获取用户信息
// 2. 检查是否已启用（防止重复启用）
// 3. 检查是否有MFA密钥（必须先调用Setup）
// 4. 验证用户提供的TOTP代码是否正确
// 5. 验证通过后，将mfa_enabled设为true
// 参数：
//   - userID: 用户ID
//   - code: TOTP验证码（6位数字）
// 返回：
//   - error: 错误信息（如已启用、密钥不存在、验证码错误）
func (s *MFAService) Enable(userID uint64, code string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrInternal
	}

	// 检查是否已启用
	if user.MFAEnabled {
		return errcode.ErrMFAAlreadyEnabled
	}

	// 检查是否有密钥
	if user.MFASecret == "" {
		return errcode.ErrMFASetupNotFound
	}

	// 验证MFA代码
	if !totp.Validate(code, user.MFASecret) {
		return errcode.ErrMFAInvalidCode
	}

	// 启用MFA
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"mfa_enabled": true,
			"updated_at":  now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Disable 禁用MFA
// 业务流程：
// 1. 获取用户信息
// 2. 检查是否已启用（未启用则无法禁用）
// 3. 验证TOTP代码（确保是本人操作，防止未授权禁用）
// 4. 禁用MFA并清除密钥（清除密钥确保安全性）
// 参数：
//   - userID: 用户ID
//   - code: TOTP验证码
// 返回：
//   - error: 错误信息
func (s *MFAService) Disable(userID uint64, code string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrInternal
	}

	// 检查是否已启用
	if !user.MFAEnabled {
		return errcode.ErrMFANotEnabled
	}

	// 验证MFA代码
	if !totp.Validate(code, user.MFASecret) {
		return errcode.ErrMFAInvalidCode
	}

	// 禁用MFA并清除密钥
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"mfa_enabled": false,
			"mfa_secret":  "",
			"updated_at":  now,
		}).Error; err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Verify 验证MFA代码
// 用于登录后的MFA二次验证
// 参数：
//   - userID: 用户ID
//   - code: TOTP验证码
// 返回：
//   - *response.MFAVerifyResp: 验证结果
//   - error: 错误信息（如MFA未启用、验证码错误）
func (s *MFAService) Verify(userID uint64, code string) (*response.MFAVerifyResp, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// 检查是否已启用
	if !user.MFAEnabled {
		return nil, errcode.ErrMFANotEnabled
	}

	// 验证MFA代码
	if !totp.Validate(code, user.MFASecret) {
		return nil, errcode.ErrMFAInvalidCode
	}

	return &response.MFAVerifyResp{
		Verified: true,
		Message:  "MFA verification successful",
	}, nil
}

// Status 获取MFA状态
// 参数：
//   - userID: 用户ID
// 返回：
//   - *response.MFAStatusResp: MFA启用状态
//   - error: 错误信息
func (s *MFAService) Status(userID uint64) (*response.MFAStatusResp, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.MFAStatusResp{
		Enabled: user.MFAEnabled,
	}, nil
}

// generateMFASecret 生成MFA密钥
func generateMFASecret() (string, error) {
	// 生成20字节的随机数
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Base32编码并移除padding
	secret := base32.StdEncoding.EncodeToString(bytes)
	secret = strings.TrimRight(secret, "=")

	return secret, nil
}
