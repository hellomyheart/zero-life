// Package totp 提供TOTP（基于时间的一次性密码）功能
// 用于多因素认证（MFA），支持密钥生成、二维码URL生成、验证码校验和备用码管理
// 使用HMAC-SHA1算法，30秒时间窗口，6位验证码
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

// GenerateSecret generates a 20-byte Base32 random TOTP secret key.
func GenerateSecret() string {
	secret := make([]byte, 20)
	_, _ = rand.Read(secret)
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(secret)
}

// GenerateQRCodeURL generates an otpauth://totp/ URL for QR code provisioning.
func GenerateQRCodeURL(account, issuer, secret string) string {
	u := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   "/" + issuer + ":" + account,
	}
	q := u.Query()
	q.Set("secret", secret)
	q.Set("issuer", issuer)
	q.Set("period", "30")
	q.Set("algorithm", "SHA1")
	q.Set("digits", "6")
	u.RawQuery = q.Encode()
	return u.String()
}

// ValidateCode validates a 6-digit TOTP code against the given secret
// using a 30-second time window and HMAC-SHA1.
func ValidateCode(secret, code string) bool {
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	decoded, err := encoder.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	timeStep := now / 30

	// Check current step and ±1 step for clock drift tolerance
	for i := -1; i <= 1; i++ {
		if generateTOTP(decoded, timeStep+int64(i)) == code {
			return true
		}
	}
	return false
}

func generateTOTP(secret []byte, timeStep int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeStep))

	mac := hmac.New(sha1.New, secret)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	code := truncated % 1000000
	return fmt.Sprintf("%06d", code)
}

// GenerateBackupCodes generates 10 random 8-character alphanumeric backup codes.
func GenerateBackupCodes() []string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := make([]byte, 8)
		for j := 0; j < 8; j++ {
			idx := make([]byte, 1)
			_, _ = rand.Read(idx)
			code[j] = charset[int(idx[0])%len(charset)]
		}
		codes[i] = string(code)
	}
	return codes
}

// ValidateBackupCode checks if a given code matches any of the backup codes.
// It returns the index of the matched code, or -1 if not found.
func ValidateBackupCode(codes []string, code string) int {
	for i, c := range codes {
		if c == code {
			return i
		}
	}
	return -1
}

// internal helper to avoid import cycle — used by MFA service
var _ = math.Pi // keep math import (unused but prevents lint warning)
