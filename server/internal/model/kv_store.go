// Package model 定义了系统的数据模型，对应数据库表结构
package model

import "time"

// KVStore 键值存储模型，对应 kv_store 表
// 用于替代 Redis，实现短期 KV 存储（带 TTL 过期机制）
// 应用场景：登录失败计数、账户锁定标记、密码重置令牌、请求限流计数
//
// 与 Redis 的对应关系：
//   - key: 对应 Redis 的键名（如 login_fail:xxx、reset_token:xxx）
//   - value: 对应 Redis 的值（如失败次数、用户ID）
//   - expires_at: 对应 Redis 的 TTL，NULL 表示永不过期
//
// 使用方式：
//   通过 KVRepository 提供的 Set/Get/Del/Incr/Exists/Expire 方法操作
//   每次操作前自动清理过期记录
type KVStore struct {
	// Key 键名，主键
	// 对应 Redis 的键名，如 "login_fail:user@example.com"
	Key string `gorm:"primaryKey;size:200" json:"key"`

	// Value 值内容
	// 对应 Redis 的值，存储为字符串
	// 如登录失败次数 "3"、密码重置令牌对应的用户ID "42"
	Value string `gorm:"not null;type:text" json:"value"`

	// ExpiresAt 过期时间
	// NULL 表示永不过期
	// 非 NULL 时，超过此时间的记录会被自动清理
	// 对应 Redis 的 TTL（EXPIRE 命令）
	ExpiresAt *time.Time `gorm:"index" json:"expires_at,omitempty"`
}

// TableName 指定 KVStore 模型对应的数据库表名为 kv_store
func (KVStore) TableName() string { return "kv_store" }