package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hellomyheart/zero-life/server/internal/config"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	UserID    uint64 `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateAccessToken(userID uint64, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(config.C.JWT.AccessTokenTTL)
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.C.JWT.Secret))
	return tokenStr, expiresAt, err
}

func (s *Service) GenerateRefreshToken(userID uint64, email string) (string, error) {
	expiresAt := time.Now().Add(config.C.JWT.RefreshTokenTTL)
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

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

func (s *Service) ParseAccessToken(tokenStr string) (*Claims, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, fmt.Errorf("invalid token type: expected access, got %s", claims.TokenType)
	}
	return claims, nil
}

func (s *Service) ParseRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type: expected refresh, got %s", claims.TokenType)
	}
	return claims, nil
}

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
