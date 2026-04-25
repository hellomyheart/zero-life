// Package middleware 提供HTTP请求中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/config"
)

// CORS 跨域资源共享中间件
// 根据配置文件中允许的Origin列表，为匹配的请求添加CORS响应头
// 参数：
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - 仅允许配置文件中CORS.AllowOrigins列表中的Origin跨域访问
//   - 允许的HTTP方法：GET, POST, PUT, PATCH, DELETE, OPTIONS
//   - 允许的请求头：Content-Type, Authorization
//   - 允许携带凭证（Cookies）
//   - 预检请求缓存时间：86400秒（1天）
//   - OPTIONS预检请求直接返回204 No Content
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false
		for _, o := range config.C.CORS.AllowOrigins {
			if o == origin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
