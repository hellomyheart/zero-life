package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
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

func (ctrl *ReportController) Tag(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Tag(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReportController) Audit(c *gin.Context) {
	userID := c.GetUint64("user_id")

	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var reconciled *bool
	if r := c.Query("reconciled"); r != "" {
		val, err := strconv.ParseBool(r)
		if err == nil {
			reconciled = &val
		}
	}

	result, err := ctrl.reportService.AuditReport(userID, accountID, startDate, endDate, reconciled)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
