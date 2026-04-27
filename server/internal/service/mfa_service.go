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
	"github.com/hellomyheart/zero-life/server/internal/pkg/hash"
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
// 验证TOTP码后启用MFA，并生成10个备用码
// 参数：
//   - userID: 用户ID
//   - code: TOTP验证码（6位数字）
// 返回：
//   - *response.BackupCodesResp: 备用码列表（明文，仅此一次展示）
//   - error: 错误信息
func (s *MFAService) Enable(userID uint64, code string) (*response.BackupCodesResp, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// 检查是否已启用
	if user.MFAEnabled {
		return nil, errcode.ErrMFAAlreadyEnabled
	}

	// 检查是否有密钥
	if user.MFASecret == "" {
		return nil, errcode.ErrMFASetupNotFound
	}

	// 验证MFA代码
	if !totp.Validate(code, user.MFASecret) {
		return nil, errcode.ErrMFAInvalidCode
	}

	// 启用MFA
	now := time.Now()
	if err := s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"mfa_enabled": true,
			"updated_at":  now,
		}).Error; err != nil {
		return nil, errcode.ErrInternal
	}

	// 生成备用码
	plainCodes, err := s.generateBackupCodes(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.BackupCodesResp{Codes: plainCodes}, nil
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

	// 验证MFA代码：先尝试TOTP，失败后尝试备用码
	if !totp.Validate(code, user.MFASecret) {
		if !s.VerifyBackupCode(userID, code) {
			return errcode.ErrMFAInvalidCode
		}
	}

	// 禁用MFA并清除密钥和备用码
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

	s.db.Where("user_id = ?", userID).Delete(&model.BackupCode{})

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

// RegenerateBackupCodes 重新生成备用码
// 先删除用户所有旧的备用码，再生成新的
// 参数：
//   - userID: 用户ID
// 返回：
//   - *response.BackupCodesResp: 新的备用码列表（明文，仅此一次展示）
//   - error: 错误信息
func (s *MFAService) RegenerateBackupCodes(userID uint64) (*response.BackupCodesResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	if !user.MFAEnabled {
		return nil, errcode.ErrMFANotEnabled
	}

	codes, err := s.generateBackupCodes(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return &response.BackupCodesResp{Codes: codes}, nil
}

// VerifyBackupCode 验证备用码
// 遍历用户所有未使用的备用码，逐个bcrypt比对
// 验证通过后标记为已使用
// 参数：
//   - userID: 用户ID
//   - code: 用户输入的备用码
// 返回：
//   - bool: 是否验证成功
func (s *MFAService) VerifyBackupCode(userID uint64, code string) bool {
	var codes []model.BackupCode
	if err := s.db.Where("user_id = ? AND used_at IS NULL", userID).Find(&codes).Error; err != nil {
		return false
	}

	for _, c := range codes {
		if hash.CheckPassword(code, c.Code) {
			now := time.Now()
			s.db.Model(&model.BackupCode{}).Where("id = ?", c.ID).Update("used_at", now)
			return true
		}
	}
	return false
}

// generateBackupCodes 为用户生成10个备用码
// 先删除旧码，再批量插入bcrypt哈希后的新码
// 返回明文备用码列表（仅此一次展示给用户）
func (s *MFAService) generateBackupCodes(userID uint64) ([]string, error) {
	s.db.Where("user_id = ?", userID).Delete(&model.BackupCode{})

	var plainCodes []string
	var records []model.BackupCode
	now := time.Now()

	for i := 0; i < 10; i++ {
		plain, hashed, err := generateOneBackupCode()
		if err != nil {
			return nil, err
		}
		plainCodes = append(plainCodes, plain)
		records = append(records, model.BackupCode{
			UserID:    userID,
			Code:      hashed,
			CreatedAt: now,
		})
	}

	if err := s.db.Create(&records).Error; err != nil {
		return nil, err
	}

	return plainCodes, nil
}

// generateOneBackupCode 生成单个备用码
// 格式：XXXX-XXXX（8位大写字母数字，中间带分隔符）
// 返回明文码、bcrypt哈希、错误
func generateOneBackupCode() (string, string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	encoded := base32.StdEncoding.EncodeToString(bytes)
	encoded = strings.TrimRight(encoded, "=")
	code := encoded[:4] + "-" + encoded[4:8]

	hashed, err := hash.HashPassword(code)
	if err != nil {
		return "", "", err
	}
	return code, hashed, nil
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
