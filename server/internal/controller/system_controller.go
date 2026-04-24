// Package controller 提供HTTP请求处理控制器
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// SystemController 系统管理控制器
// 处理系统管理相关的HTTP请求
type SystemController struct {
	systemService *service.SystemService
}

// NewSystemController 创建系统管理控制器实例
// 参数：
//   - systemService: 系统管理业务服务
// 返回：
//   - *SystemController: 系统管理控制器实例
func NewSystemController(systemService *service.SystemService) *SystemController {
	return &SystemController{systemService: systemService}
}

// GetSystemInfo 获取系统信息
// 参数：
//   - c: Gin上下文
func (ctrl *SystemController) GetSystemInfo(c *gin.Context) {
	result, err := ctrl.systemService.GetSystemInfo()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// HealthCheck 健康检查
// 参数：
//   - c: Gin上下文
func (ctrl *SystemController) HealthCheck(c *gin.Context) {
	result, err := ctrl.systemService.HealthCheck()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
