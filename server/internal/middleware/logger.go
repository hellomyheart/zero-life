// Package middleware 提供HTTP请求中间件
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 请求日志中间件
// 记录每个HTTP请求的方法、路径、状态码、响应时间和客户端IP
// 参数：
//   - logger: Zap日志实例，用于输出结构化日志
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - 在请求处理前记录开始时间
//   - 在请求处理完成后计算响应延迟并输出日志
//   - 日志级别为Info
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
