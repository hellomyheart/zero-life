package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

const cronToken = "zero-life-cron-2024"

type CronController struct {
	cronService *service.CronService
}

func NewCronController(cronService *service.CronService) *CronController {
	return &CronController{cronService: cronService}
}

func (ctrl *CronController) Run(c *gin.Context) {
	token := c.Param("token")
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
