// Package model 定义了系统的数据模型
// 包含用户、账户、交易、分类、标签等核心实体
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// 存储用户基本信息、偏好设置和安全配置
type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`                    // 用户ID，主键自增
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`            // 邮箱地址，唯一索引
	Password  string         `gorm:"not null;size:255" json:"-"`                            // 密码（加密存储），不返回给前端
	Nickname  string         `gorm:"not null;size:100" json:"nickname"`                     // 用户昵称
	Language  string         `gorm:"size:10;default:zh-CN" json:"language"`                 // 语言偏好，默认中文
	Timezone  string         `gorm:"size:50;default:Asia/Shanghai" json:"timezone"`         // 时区设置，默认上海时区
	Role      string         `gorm:"size:20;default:user" json:"role"`                      // 用户角色（user/admin）
	MFASecret string         `gorm:"size:255" json:"-"`                                     // MFA密钥，不返回给前端
	MFAEnabled bool          `gorm:"default:false" json:"mfa_enabled"`                      // 是否启用多因素认证
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`                            // 创建时间
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`                            // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                        // 软删除时间
}
