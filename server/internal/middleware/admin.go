// Package middleware 提供HTTP请求中间件
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// Admin 管理员权限校验中间件
// 检查当前认证用户是否具有admin角色，非管理员返回403
// 参数：
//   - db: GORM数据库实例，用于查询用户角色信息
// 返回：
//   - gin.HandlerFunc: Gin中间件函数
// 业务规则：
//   - 必须先经过Auth中间件，确保上下文中存在user_id
//   - 从数据库查询用户信息，验证角色为admin
//   - 非管理员请求将被中止并返回403 Forbidden
func Admin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
			})
			return
		}

		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "user not found",
			})
			return
		}

		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "admin access required",
			})
			return
		}

		c.Set("role", user.Role)
		c.Next()
	}
}
