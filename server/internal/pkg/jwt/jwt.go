package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zero-life/server/internal/config"
)

type Claims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
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
	return s.parseToken(tokenStr)
}

func (s *Service) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return s.parseToken(tokenStr)
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
