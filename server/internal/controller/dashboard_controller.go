// Package controller 提供HTTP请求处理控制器
// DashboardController 仪表盘控制器
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// DashboardController 仪表盘控制器
// 处理仪表盘数据查询请求，包括月度收支、总余额、预算预警、账单提醒和最近交易
type DashboardController struct {
	dashboardService *service.DashboardService // 仪表盘业务服务
}

// NewDashboardController 创建仪表盘控制器实例
// 参数：
//   - dashboardService: 仪表盘业务服务实例
// 返回：
//   - *DashboardController: 仪表盘控制器实例
func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

// Get 获取仪表盘数据
// 返回当前用户的仪表盘汇总数据，包括月度收支、总余额、预算预警、账单提醒和最近交易
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：DashboardResp（含月度收支、总余额、预算预警、账单提醒、最近交易）
// @Summary      Get dashboard data
// @Description  Get dashboard summary including monthly income/expense, total balance, budget alerts, bill reminders and recent transactions
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/dashboard [get]
// @Security     BearerAuth
func (ctrl *DashboardController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.dashboardService.Get(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
