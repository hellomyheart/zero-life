// Package middleware 提供HTTP请求中间件
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
)

// RateLimit 基于Redis的请求限流中间件
// 使用Redis计数器实现滑动窗口限流，按客户端IP进行限制
// 参数：
//   - rdb: Redis客户端实例，用于存储请求计数
//   - limit: 时间窗口内允许的最大请求数
//   - window: 限流时间窗口
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - 按客户端IP地址进行限流，key格式为"ratelimit:<IP>"
//   - 首次请求时设置计数器过期时间为window
//   - 超过limit的请求返回429 Too Many Requests
//   - Redis操作失败时放行请求（降级策略）
func RateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()
		ctx := context.Background()

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": errcode.ErrTooManyReqs.Code, "message": errcode.ErrTooManyReqs.Message})
			c.Abort()
			return
		}

		c.Next()
	}
}
