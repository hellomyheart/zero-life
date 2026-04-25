// Package hash 提供密码哈希和验证功能
// 使用bcrypt算法进行密码加密，cost因子为10
package hash

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行bcrypt哈希加密
// 参数：
//   - password: 明文密码
// 返回：
//   - string: bcrypt哈希后的密码字符串
//   - error: 哈希失败时返回错误
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) // cost=10，平衡安全性和性能
	return string(bytes), err
}

// CheckPassword 验证明文密码与哈希密码是否匹配
// 参数：
//   - password: 待验证的明文密码
//   - hash: 存储的bcrypt哈希密码
// 返回：
//   - bool: 匹配返回true，不匹配返回false
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
