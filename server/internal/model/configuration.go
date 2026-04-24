// Package model 数据模型层，对应数据库结构
package model

import "time"

// Configuration 系统配置模型，对应configurations表
// 存储全局系统配置项，采用键值对（Name-Value）结构
// 例如：系统名称、默认货币、注册开关等全局设置
// 与Preference（用户偏好）不同，Configuration是系统级别的配置
type Configuration struct {
	// ID 配置项唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Name 配置项名称（键），如"site_name"、"default_currency"
	// gorm:"uniqueIndex" 创建唯一索引，确保配置项名称不重复
	// gorm:"size:100" 限制名称最大长度为100个字符
	Name string `gorm:"uniqueIndex;not null;size:100" json:"name"`

	// Value 配置项的值（值），存储配置的具体内容
	// gorm:"type:text" 指定数据库字段类型为TEXT，可存储较长的文本内容
	// 与VARCHAR不同，TEXT类型不限制长度，适合存储JSON等长文本
	Value string `gorm:"type:text" json:"value"`

	// CreatedAt 配置项创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 配置项最后更新时间，每次修改配置时自动更新
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定数据库表名为configurations
func (Configuration) TableName() string { return "configurations" }
