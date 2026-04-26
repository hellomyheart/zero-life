// Package middleware 提供HTTP请求中间件
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// RateLimit 基于键值存储的请求限流中间件
// 使用 kv_store 表实现滑动窗口限流，按客户端IP进行限制
// 参数：
//   - kvRepo: 键值存储仓库实例，用于存储请求计数
//   - limit: 时间窗口内允许的最大请求数
//   - window: 限流时间窗口
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - 按客户端IP地址进行限流，key格式为"ratelimit:<IP>"
//   - 首次请求时设置计数器过期时间为window
//   - 超过limit的请求返回429 Too Many Requests
//   - 存储操作失败时放行请求（降级策略）
func RateLimit(kvRepo *repository.KVRepository, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()

		count, err := kvRepo.Incr(key)
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			kvRepo.Expire(key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": errcode.ErrTooManyReqs.Code, "message": errcode.ErrTooManyReqs.Message})
			c.Abort()
			return
		}

		c.Next()
	}
}