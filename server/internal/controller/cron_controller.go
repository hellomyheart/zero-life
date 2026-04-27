package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type CronController struct {
	cronService *service.CronService
}

func NewCronController(cronService *service.CronService) *CronController {
	return &CronController{cronService: cronService}
}

func (c *CronController) List(ctx *gin.Context) {
	tasks := c.cronService.ListTasks()
	Success(ctx, tasks)
}

func (c *CronController) Run(ctx *gin.Context) {
	taskID := ctx.Param("id")
	result, err := c.cronService.RunTask(taskID)
	if err != nil {
		Error(ctx, http.StatusNotFound, errcode.ErrNotFound)
		return
	}
	Success(ctx, result)
}
