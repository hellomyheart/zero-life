package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/service"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (ctrl *ReportController) IncomeExpense(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.IncomeExpense(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReportController) Category(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Category(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReportController) Budget(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Budget(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReportController) NetWorth(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.NetWorth(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReportController) Trend(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Trend(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
