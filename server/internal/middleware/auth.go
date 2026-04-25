// Package middleware 提供HTTP请求中间件
// 包含认证、权限校验、跨域、日志、限流等中间件
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
)

// Auth JWT认证中间件
// 从请求头中提取Bearer Token，解析JWT令牌，将用户ID和邮箱注入到Gin上下文中
// 参数：
//   - jwtService: JWT服务实例，用于解析和验证令牌
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - Authorization头必须存在且格式为"Bearer <token>"
//   - 令牌必须有效且未过期
//   - 解析成功后将user_id和email注入上下文，供后续处理器使用
func Auth(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": errcode.ErrUnauthorized.Code, "message": errcode.ErrUnauthorized.Message})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": errcode.ErrUnauthorized.Code, "message": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := jwtService.ParseAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": errcode.ErrUnauthorized.Code, "message": errcode.ErrInvalidToken.Message})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
