package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/pkg/hash"
	"github.com/zero-life/server/internal/pkg/jwt"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

const (
	loginFailKeyPrefix  = "login_fail:"
	loginLockKeyPrefix  = "login_lock:"
	resetTokenKeyPrefix = "reset_token:"
	maxLoginFails       = 5
	loginLockDuration   = 30 * time.Minute
	resetTokenTTL       = 24 * time.Hour
)

type AuthService struct {
	authRepo   *repository.AuthRepository
	jwtService *jwt.Service
	rdb        *redis.Client
}

func NewAuthService(authRepo *repository.AuthRepository, jwtService *jwt.Service, rdb *redis.Client) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtService: jwtService,
		rdb:        rdb,
	}
}

func (s *AuthService) Register(req *request.RegisterReq) (*response.LoginResp, error) {
	// Check email uniqueness
	existing, err := s.authRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errcode.ErrEmailExists
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errcode.ErrInternal
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Create user
	user := &model.User{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
	}
	if err := s.authRepo.Create(user); err != nil {
		return nil, errcode.ErrInternal
	}

	// Generate token pair
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

func (s *AuthService) Login(req *request.LoginReq) (*response.LoginResp, error) {
	ctx := context.Background()

	// Check if account is locked
	lockKey := loginLockKeyPrefix + req.Email
	locked, err := s.rdb.Exists(ctx, lockKey).Result()
	if err == nil && locked > 0 {
		return nil, errcode.ErrAccountLocked
	}

	// Find user
	user, err := s.authRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrInvalidCredential
		}
		return nil, errcode.ErrInternal
	}

	// Check password
	if !hash.CheckPassword(req.Password, user.Password) {
		// Record failed attempt
		failKey := loginFailKeyPrefix + req.Email
		count, _ := s.rdb.Incr(ctx, failKey).Result()
		if count == 1 {
			s.rdb.Expire(ctx, failKey, loginLockDuration)
		}
		if count >= maxLoginFails {
			s.rdb.Set(ctx, lockKey, "1", loginLockDuration)
			s.rdb.Del(ctx, failKey)
			return nil, errcode.ErrAccountLocked
		}
		return nil, errcode.ErrInvalidCredential
	}

	// Clear fail count on success
	failKey := loginFailKeyPrefix + req.Email
	s.rdb.Del(ctx, failKey)

	// Generate token pair
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

func (s *AuthService) RefreshToken(refreshToken string) (*response.LoginResp, error) {
	claims, err := s.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}

	// Verify user still exists
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

func (s *AuthService) ForgotPassword(email string) error {
	// Check user exists
	user, err := s.authRepo.FindByEmail(email)
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

	// Store in Redis with 24h TTL
	ctx := context.Background()
	key := resetTokenKeyPrefix + token
	if err := s.rdb.Set(ctx, key, fmt.Sprintf("%d", user.ID), resetTokenTTL).Err(); err != nil {
		return errcode.ErrInternal
	}

	// TODO: Send email with reset link
	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	ctx := context.Background()
	key := resetTokenKeyPrefix + token

	// Verify token
	userIDStr, err := s.rdb.Get(ctx, key).Result()
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
	s.rdb.Del(ctx, key)

	return nil
}

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
