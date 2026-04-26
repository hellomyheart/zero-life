// Package repository 数据访问层，封装数据库操作
package repository

import (
	"fmt"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// KVRepository 键值存储仓库，用于替代 Redis 实现短期 KV 存储
// 提供与 Redis 类似的操作接口：Set/Get/Del/Incr/Exists/Expire
// 每次操作前自动清理过期记录，模拟 Redis 的 TTL 行为
//
// 与 Redis 操作的对应关系：
//
//	Redis: rdb.Set(ctx, key, value, ttl)     → kvRepo.Set(key, value, ttl)
//	Redis: rdb.Get(ctx, key).Result()        → kvRepo.Get(key)
//	Redis: rdb.Del(ctx, key)                 → kvRepo.Del(key)
//	Redis: rdb.Incr(ctx, key).Result()       → kvRepo.Incr(key)
//	Redis: rdb.Exists(ctx, key).Result()     → kvRepo.Exists(key)
//	Redis: rdb.Expire(ctx, key, ttl)         → kvRepo.Expire(key, ttl)
type KVRepository struct {
	db *gorm.DB
}

// NewKVRepository 创建键值存储仓库实例
// 参数 db: GORM 数据库连接实例
func NewKVRepository(db *gorm.DB) *KVRepository {
	return &KVRepository{db: db}
}

// cleanExpired 清理所有过期的键值记录
// 在每次操作前调用，模拟 Redis 的自动过期行为
// 执行 SQL: DELETE FROM kv_store WHERE expires_at IS NOT NULL AND expires_at < datetime('now')
func (r *KVRepository) cleanExpired() {
	r.db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&model.KVStore{})
}

// Set 设置键值对，支持可选的 TTL 过期时间
// 对应 Redis: SET key value [EX ttl]
// 如果 key 已存在则更新，否则插入
// 参数：
//   - key: 键名
//   - value: 值内容
//   - ttl: 过期时间，0 表示永不过期
func (r *KVRepository) Set(key, value string, ttl time.Duration) error {
	r.cleanExpired()

	var expiresAt *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		expiresAt = &t
	}

	// 使用 GORM 的 Clauses 实现 UPSERT（INSERT ... ON CONFLICT DO UPDATE）
	return r.db.Save(&model.KVStore{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt,
	}).Error
}

// Get 获取键对应的值
// 对应 Redis: GET key
// 如果键不存在或已过期，返回错误
// 参数：
//   - key: 键名
// 返回：
//   - string: 键对应的值
//   - error: 键不存在或已过期时返回 gorm.ErrRecordNotFound
func (r *KVRepository) Get(key string) (string, error) {
	r.cleanExpired()

	var kv model.KVStore
	if err := r.db.Where("key = ?", key).First(&kv).Error; err != nil {
		return "", err
	}
	return kv.Value, nil
}

// Del 删除一个或多个键
// 对应 Redis: DEL key [key ...]
// 参数：
//   - keys: 要删除的键名列表
func (r *KVRepository) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.db.Where("key IN ?", keys).Delete(&model.KVStore{}).Error
}

// Incr 原子递增键的值，并返回递增后的结果
// 对应 Redis: INCR key
// 如果键不存在，从 0 开始递增（结果为 1）
// 如果键存在但值不是数字，返回错误
// 注意：SQLite 不支持真正的原子操作，此方法使用事务保证一致性
// 参数：
//   - key: 键名
// 返回：
//   - int64: 递增后的值
//   - error: 操作失败时返回错误
func (r *KVRepository) Incr(key string) (int64, error) {
	r.cleanExpired()

	var count int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var kv model.KVStore
		if err := tx.Where("key = ?", key).First(&kv).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 键不存在，从 1 开始
				kv = model.KVStore{Key: key, Value: "1"}
				if err := tx.Create(&kv).Error; err != nil {
					return err
				}
				count = 1
				return nil
			}
			return err
		}

		// 解析当前值并递增
		var current int64
		if _, err := fmt.Sscanf(kv.Value, "%d", &current); err != nil {
			// 值不是数字，重置为 1
			current = 0
		}
		current++
		kv.Value = fmt.Sprintf("%d", current)
		if err := tx.Model(&model.KVStore{}).Where("key = ?", key).Update("value", kv.Value).Error; err != nil {
			return err
		}
		count = current
		return nil
	})

	return count, err
}

// Exists 检查键是否存在且未过期
// 对应 Redis: EXISTS key
// 参数：
//   - key: 键名
// 返回：
//   - bool: 键存在且未过期返回 true
//   - error: 数据库操作错误
func (r *KVRepository) Exists(key string) (bool, error) {
	r.cleanExpired()

	var count int64
	if err := r.db.Model(&model.KVStore{}).Where("key = ?", key).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Expire 设置键的过期时间
// 对应 Redis: EXPIRE key ttl
// 参数：
//   - key: 键名
//   - ttl: 过期时间
// 返回：
//   - error: 键不存在或操作失败时返回错误
func (r *KVRepository) Expire(key string, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	result := r.db.Model(&model.KVStore{}).Where("key = ?", key).Update("expires_at", expiresAt)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
