// Package model 数据模型层，对应数据库结构
package model

import "time"

// BackupCode MFA备用码模型，对应backup_codes表
// 存储两步验证（MFA）的备用恢复码，当用户无法使用认证器时可用备用码登录
// 每个备用码只能使用一次，使用后标记UsedAt时间
type BackupCode struct {
	// ID 备用码唯一标识，主键自增
	// gorm:"primaryKey" 表示这是主键
	// gorm:"autoIncrement" 表示主键自动递增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户ID，关联users表
	// gorm:"not null" 表示数据库字段不允许为空
	// gorm:"index" 为此字段创建索引，加速按用户查询备用码
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Code 备用码字符串，通常是随机生成的8-10位字母数字组合
	// gorm:"size:20" 限制数据库字段最大长度为20个字符
	// gorm:"index" 为此字段创建索引，登录验证时快速查找备用码
	Code string `gorm:"not null;size:20;index" json:"code"`

	// UsedAt 备用码使用时间，指针类型表示可为nil（未使用）
	// nil表示备用码尚未使用，非nil表示已使用及使用时间
	UsedAt *time.Time `json:"used_at"`

	// CreatedAt 备用码创建时间，即生成备用码的时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定数据库表名为backup_codes
func (BackupCode) TableName() string { return "backup_codes" }
