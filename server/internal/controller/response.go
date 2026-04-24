// Package controller 提供HTTP请求处理控制器
// 响应工具函数，统一封装HTTP成功和错误响应格式
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": data})
}

func SuccessPage(c *gin.Context, result *pagination.Result) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

func Error(c *gin.Context, httpStatus int, err *errcode.Error) {
	c.JSON(httpStatus, gin.H{"code": err.Code, "message": err.Message})
}
