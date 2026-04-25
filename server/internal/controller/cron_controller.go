// Package controller 提供HTTP请求处理控制器
// CronController 定时任务控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/config"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// CronController 定时任务控制器
// 处理外部定时任务触发请求（如cron定时调用），需要验证令牌防止未授权访问
type CronController struct {
	cronService *service.CronService
}

// NewCronController 创建定时任务控制器实例
func NewCronController(cronService *service.CronService) *CronController {
	return &CronController{cronService: cronService}
}

// Run 执行定时任务
// 通过URL路径中的token参数验证请求合法性，防止未授权触发
// token从配置文件app.cron_token读取，不再硬编码
func (ctrl *CronController) Run(c *gin.Context) {
	token := c.Param("token")
	// 从配置读取cron令牌，如果未配置则使用默认值（建议在配置文件中设置）
	cronToken := config.C.App.CronToken
	if cronToken == "" {
		cronToken = "zero-life-cron-2024"
	}
	if token != cronToken {
		Error(c, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

	result, err := ctrl.cronService.CronRun()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
