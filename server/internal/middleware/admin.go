package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// Admin checks if the authenticated user has admin role.
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

		c.Next()
	}
}
