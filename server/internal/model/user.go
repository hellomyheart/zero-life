// Package model 定义了系统的数据模型
// 包含用户、账户、交易、分类、标签等核心实体
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型，对应 users 表
// 
// 功能说明：
// - 存储用户基本信息（邮箱、昵称）
// - 存储用户偏好设置（语言、时区）
// - 存储安全配置（密码、MFA）
// - 支持多角色（普通用户/管理员）
//
// 用户角色：
// - user: 普通用户，只能管理自己的数据
// - admin: 管理员，可以管理所有用户和系统配置
//
// 安全特性：
// - 密码使用 bcrypt 加密存储
// - 支持多因素认证（MFA）
// - 软删除保护用户数据
//
// 使用场景：
// - 用户认证（登录/注册）
// - 个性化设置（语言/时区）
// - 权限控制（角色管理）
// - 数据隔离（UserID 关联）
type User struct {
	// ID 用户唯一标识，主键自增
	// 所有用户相关数据都通过 UserID 关联到此字段
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// Email 用户邮箱地址
	// 用于登录和接收系统通知
	// gorm:"uniqueIndex" 创建唯一索引，确保邮箱不重复
	// gorm:"size:255" 限制最大长度为 255 个字符
	Email string `gorm:"uniqueIndex;not null;size:255" json:"email"`
	
	// Password 用户密码（加密存储）
	// 使用 bcrypt 算法加密，确保安全性
	// gorm:"size:255" 存储 bcrypt 哈希值（60 字符）+ 余量
	// json:"-" 不返回给前端，保护敏感信息
	Password string `gorm:"not null;size:255" json:"-"`
	
	// Nickname 用户昵称
	// 显示在界面上的友好名称
	// gorm:"size:100" 限制最大长度为 100 个字符
	Nickname string `gorm:"not null;size:100" json:"nickname"`
	
	// Language 用户语言偏好
	// 用于前端界面国际化
	// 可选值：zh-CN（中文）, en-US（英文）
	// 默认值：zh-CN
	// gorm:"size:10" 限制最大长度为 10 个字符
	Language string `gorm:"size:10;default:zh-CN" json:"language"`
	
	// Timezone 用户时区设置
	// 用于日期时间显示和计算
	// 可选值：Asia/Shanghai（中国）, America/New_York（美东）等
	// 默认值：Asia/Shanghai
	// gorm:"size:50" 限制最大长度为 50 个字符
	Timezone string `gorm:"size:50;default:Asia/Shanghai" json:"timezone"`
	
	// Role 用户角色
	// 控制用户权限级别
	// 可选值：
	//   - user: 普通用户，只能管理自己的数据
	//   - admin: 管理员，可以管理所有用户和系统配置
	// 默认值：user
	// gorm:"size:20" 限制最大长度为 20 个字符
	Role string `gorm:"size:20;default:user" json:"role"`
	
	// MFASecret MFA（多因素认证）密钥
	// 用于生成 TOTP 验证码
	// 使用 Google Authenticator 等验证器 App
	// gorm:"size:255" 存储 Base32 编码的密钥
	// json:"-" 不返回给前端，保护敏感信息
	MFASecret string `gorm:"size:255" json:"-"`
	
	// MFAEnabled 是否启用多因素认证
	// true: 登录时需要输入 6 位验证码
	// false: 仅需密码即可登录
	// 默认值：false（不启用）
	MFAEnabled bool `gorm:"default:false" json:"mfa_enabled"`
	
	// CreatedAt 用户创建时间
	// 记录用户注册时间，不可变
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 用户最后更新时间
	// 记录用户信息最后修改时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 User 模型对应的数据库表名为 users
func (User) TableName() string { return "users" }
