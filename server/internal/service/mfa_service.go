package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
	"github.com/hellomyheart/zero-life/server/internal/pkg/totp"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

const (
	mfaSetupKeyPrefix  = "mfa_setup:"
	mfaTokenKeyPrefix  = "mfa_token:"
	mfaSetupTTL        = 5 * time.Minute
	mfaTokenTTL        = 5 * time.Minute
)

type MFAService struct {
	authRepo      *repository.AuthRepository
	backupCodeRepo *repository.BackupCodeRepository
	jwtService    *jwt.Service
	rdb           *redis.Client
}

func NewMFAService(
	authRepo *repository.AuthRepository,
	backupCodeRepo *repository.BackupCodeRepository,
	jwtService *jwt.Service,
	rdb *redis.Client,
) *MFAService {
	return &MFAService{
		authRepo:      authRepo,
		backupCodeRepo: backupCodeRepo,
		jwtService:    jwtService,
		rdb:           rdb,
	}
}

// Enable generates a TOTP secret and backup codes, stores them in Redis temporarily.
func (s *MFAService) Enable(userID uint64) (*response.MFAEnableResp, error) {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if user.MFAEnabled {
		return nil, errcode.ErrMFAAlreadyEnabled
	}

	secret := totp.GenerateSecret()
	qrURL := totp.GenerateQRCodeURL(user.Email, "ZeroLife", secret)
	backupCodes := totp.GenerateBackupCodes()

	// Store setup data in Redis temporarily
	ctx := context.Background()
	key := mfaSetupKeyPrefix + fmt.Sprintf("%d", userID)
	setupData := fmt.Sprintf("%s", secret)
	if err := s.rdb.Set(ctx, key, setupData, mfaSetupTTL).Err(); err != nil {
		return nil, errcode.ErrInternal
	}

	// Store backup codes in Redis temporarily
	backupKey := mfaSetupKeyPrefix + "backup:" + fmt.Sprintf("%d", userID)
	backupData := ""
	for i, code := range backupCodes {
		if i > 0 {
			backupData += ","
		}
		backupData += code
	}
	if err := s.rdb.Set(ctx, backupKey, backupData, mfaSetupTTL).Err(); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.MFAEnableResp{
		Secret:      secret,
		QRCodeURL:   qrURL,
		BackupCodes: backupCodes,
	}, nil
}

// Confirm verifies the TOTP code and officially enables MFA for the user.
func (s *MFAService) Confirm(userID uint64, code string) error {
	ctx := context.Background()

	// Get pending setup data from Redis
	key := mfaSetupKeyPrefix + fmt.Sprintf("%d", userID)
	secret, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return errcode.ErrMFASetupNotFound
	}

	// Validate the TOTP code
	if !totp.ValidateCode(secret, code) {
		return errcode.ErrMFAInvalidCode
	}

	// Get backup codes from Redis
	backupKey := mfaSetupKeyPrefix + "backup:" + fmt.Sprintf("%d", userID)
	backupData, err := s.rdb.Get(ctx, backupKey).Result()
	if err != nil {
		return errcode.ErrMFASetupNotFound
	}

	// Update user: enable MFA and store secret
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrNotFound
	}
	user.MFAEnabled = true
	user.MFASecret = secret
	if err := s.authRepo.Update(user); err != nil {
		return errcode.ErrInternal
	}

	// Save backup codes to database
	codes := splitBackupCodes(backupData)
	backupCodes := make([]model.BackupCode, 0, len(codes))
	for _, c := range codes {
		backupCodes = append(backupCodes, model.BackupCode{
			UserID: userID,
			Code:   c,
		})
	}
	if err := s.backupCodeRepo.Create(backupCodes); err != nil {
		return errcode.ErrInternal
	}

	// Clean up Redis keys
	s.rdb.Del(ctx, key)
	s.rdb.Del(ctx, backupKey)

	return nil
}

// Disable verifies the TOTP code or backup code, then disables MFA.
func (s *MFAService) Disable(userID uint64, code string) error {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if !user.MFAEnabled {
		return errcode.ErrMFANotEnabled
	}

	if !s.verifyCodeOrBackup(userID, user.MFASecret, code) {
		return errcode.ErrMFAInvalidCode
	}

	// Disable MFA
	user.MFAEnabled = false
	user.MFASecret = ""
	if err := s.authRepo.Update(user); err != nil {
		return errcode.ErrInternal
	}

	// Delete all backup codes
	if err := s.backupCodeRepo.DeleteByUserID(userID); err != nil {
		return errcode.ErrInternal
	}

	return nil
}

// Verify validates a TOTP code or backup code (used during login).
func (s *MFAService) Verify(userID uint64, code string) bool {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return false
	}
	if !user.MFAEnabled {
		return false
	}
	return s.verifyCodeOrBackup(userID, user.MFASecret, code)
}

// VerifyMFA verifies the MFA code after login and issues a real access token.
func (s *MFAService) VerifyMFA(mfaToken, code string) (*response.LoginResp, error) {
	ctx := context.Background()

	// Get user ID from MFA token in Redis
	key := mfaTokenKeyPrefix + mfaToken
	userIDStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, errcode.ErrMFAInvalidToken
	}

	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)

	// Verify the code
	if !s.Verify(userID, code) {
		return nil, errcode.ErrMFAInvalidCode
	}

	// Get user and generate real token pair
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Delete used MFA token
	s.rdb.Del(ctx, key)

	return &response.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// GetBackupCodes returns the list of backup codes for the user.
func (s *MFAService) GetBackupCodes(userID uint64) ([]response.MFABackupCodeResp, error) {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if !user.MFAEnabled {
		return nil, errcode.ErrMFANotEnabled
	}

	codes, err := s.backupCodeRepo.ListByUserID(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.MFABackupCodeResp, 0, len(codes))
	for _, c := range codes {
		var usedAt *string
		if c.UsedAt != nil {
			s := c.UsedAt.Format("2006-01-02 15:04:05")
			usedAt = &s
		}
		items = append(items, response.MFABackupCodeResp{
			ID:        c.ID,
			Code:      c.Code,
			UsedAt:    usedAt,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

// RegenerateBackupCodes verifies the code and generates new backup codes.
func (s *MFAService) RegenerateBackupCodes(userID uint64, code string) ([]response.MFABackupCodeResp, error) {
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if !user.MFAEnabled {
		return nil, errcode.ErrMFANotEnabled
	}

	if !s.verifyCodeOrBackup(userID, user.MFASecret, code) {
		return nil, errcode.ErrMFAInvalidCode
	}

	// Delete old backup codes
	if err := s.backupCodeRepo.DeleteByUserID(userID); err != nil {
		return nil, errcode.ErrInternal
	}

	// Generate new backup codes
	newCodes := totp.GenerateBackupCodes()
	backupCodes := make([]model.BackupCode, 0, len(newCodes))
	for _, c := range newCodes {
		backupCodes = append(backupCodes, model.BackupCode{
			UserID: userID,
			Code:   c,
		})
	}
	if err := s.backupCodeRepo.Create(backupCodes); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.GetBackupCodes(userID)
}

// StoreMFAToken stores a temporary MFA token in Redis after login when MFA is required.
func (s *MFAService) StoreMFAToken(userID uint64) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", errcode.ErrInternal
	}
	mfaToken := hex.EncodeToString(tokenBytes)

	ctx := context.Background()
	key := mfaTokenKeyPrefix + mfaToken
	if err := s.rdb.Set(ctx, key, fmt.Sprintf("%d", userID), mfaTokenTTL).Err(); err != nil {
		return "", errcode.ErrInternal
	}

	return mfaToken, nil
}

func (s *MFAService) verifyCodeOrBackup(userID uint64, secret, code string) bool {
	// Try TOTP code first
	if totp.ValidateCode(secret, code) {
		return true
	}

	// Try backup code
	bc, err := s.backupCodeRepo.FindByCode(userID, code)
	if err != nil {
		return false
	}

	// Mark backup code as used
	if err := s.backupCodeRepo.MarkUsed(bc.ID); err != nil {
		return false
	}

	return true
}

func splitBackupCodes(data string) []string {
	if data == "" {
		return nil
	}
	result := make([]string, 0)
	current := ""
	for _, ch := range data {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
