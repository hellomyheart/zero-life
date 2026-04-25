// Package jwt 提供JWT令牌的生成和解析功能
// 支持访问令牌（AccessToken）和刷新令牌（RefreshToken）的生成与验证
// 使用HMAC-SHA256签名算法，密钥从配置文件读取
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hellomyheart/zero-life/server/internal/config"
)

// Claims JWT令牌的自定义声明结构
// 包含用户ID和邮箱信息，以及标准的JWT注册声明
type Claims struct {
	UserID uint64 `json:"user_id"` // 用户ID，用于标识用户身份
	Email  string `json:"email"`   // 用户邮箱，用于显示和通知
	jwt.RegisteredClaims
}

// TokenPair 令牌对结构
// 包含访问令牌、刷新令牌和访问令牌的过期时间
type TokenPair struct {
	AccessToken  string    `json:"access_token"`  // 访问令牌，用于API请求认证
	RefreshToken string    `json:"refresh_token"` // 刷新令牌，用于获取新的访问令牌
	ExpiresAt    time.Time `json:"expires_at"`    // 访问令牌的过期时间
}

// Service JWT服务
// 提供令牌的生成和解析功能
type Service struct{}

// NewService 创建JWT服务实例
// 返回：
//   - *Service: JWT服务实例
func NewService() *Service {
	return &Service{}
}

// GenerateAccessToken 生成访问令牌
// 根据用户ID和邮箱生成短期访问令牌
// 参数：
//   - userID: 用户ID
//   - email: 用户邮箱
// 返回：
//   - string: 生成的访问令牌字符串
//   - time.Time: 令牌过期时间
//   - error: 生成失败时返回错误
func (s *Service) GenerateAccessToken(userID uint64, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(config.C.JWT.AccessTokenTTL)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.C.JWT.Secret))
	return tokenStr, expiresAt, err
}

// GenerateRefreshToken 生成刷新令牌
// 根据用户ID和邮箱生成长期刷新令牌，用于获取新的访问令牌
// 参数：
//   - userID: 用户ID
//   - email: 用户邮箱
// 返回：
//   - string: 生成的刷新令牌字符串
//   - error: 生成失败时返回错误
func (s *Service) GenerateRefreshToken(userID uint64, email string) (string, error) {
	expiresAt := time.Now().Add(config.C.JWT.RefreshTokenTTL)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// GenerateTokenPair 生成令牌对
// 同时生成访问令牌和刷新令牌
// 参数：
//   - userID: 用户ID
//   - email: 用户邮箱
// 返回：
//   - *TokenPair: 令牌对（含访问令牌、刷新令牌和过期时间）
//   - error: 生成失败时返回错误
func (s *Service) GenerateTokenPair(userID uint64, email string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.GenerateRefreshToken(userID, email)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ParseAccessToken 解析访问令牌
// 参数：
//   - tokenStr: 访问令牌字符串
// 返回：
//   - *Claims: 解析后的令牌声明
//   - error: 解析失败时返回错误
func (s *Service) ParseAccessToken(tokenStr string) (*Claims, error) {
	return s.parseToken(tokenStr)
}

// ParseRefreshToken 解析刷新令牌
// 参数：
//   - tokenStr: 刷新令牌字符串
// 返回：
//   - *Claims: 解析后的令牌声明
//   - error: 解析失败时返回错误
func (s *Service) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return s.parseToken(tokenStr)
}

// parseToken 解析JWT令牌的内部实现
// 使用配置文件中的密钥验证令牌签名和有效期
// 参数：
//   - tokenStr: 令牌字符串
// 返回：
//   - *Claims: 解析后的令牌声明
//   - error: 令牌无效或过期时返回错误
func (s *Service) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
