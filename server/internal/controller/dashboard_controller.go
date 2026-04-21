package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zero-life/server/internal/service"
)

type DashboardController struct {
	dashboardService *service.DashboardService
}

func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

func (ctrl *DashboardController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.dashboardService.Get(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
