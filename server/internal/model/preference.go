// Package model 数据模型层，对应数据库结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// Preference 用户偏好模型，对应preferences表
// 存储用户个性化配置，采用键值对（Key-Value）结构
// 与Configuration（系统配置）不同，Preference是用户级别的配置
// 例如：用户默认账户、列表每页条数、日期格式等个人偏好
type Preference struct {
	// ID 偏好记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户ID
	// gorm:"uniqueIndex:idx_pref_user_key" 创建复合唯一索引的一部分
	// 与Key字段组合，确保同一用户不会有重复的偏好键
	UserID uint64 `gorm:"not null;uniqueIndex:idx_pref_user_key" json:"user_id"`

	// Key 偏好键名，如"list_page_size"、"default_account"等
	// gorm:"uniqueIndex:idx_pref_user_key" 与UserID组成复合唯一索引
	// 复合唯一索引保证：同一用户的同一偏好键只能存在一条记录
	// gorm:"size:100" 限制键名最大长度为100个字符
	Key string `gorm:"not null;size:100;uniqueIndex:idx_pref_user_key" json:"key"`

	// Value 偏好值，存储配置的具体内容
	// gorm:"type:text" 指定数据库字段类型为TEXT，可存储较长的文本
	// 例如JSON格式的复杂配置值
	Value string `gorm:"type:text" json:"value"`

	// CreatedAt 偏好创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 偏好最后更新时间，用户修改偏好时自动更新
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间，GORM软删除字段
	// gorm.DeletedAt 记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM查询时自动过滤已软删除的记录
	// json:"-" JSON序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名为preferences
func (Preference) TableName() string { return "preferences" }
