// Package controller 控制器层
package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// InsightController 数据洞察控制器
// 提供深度数据分析功能，包括支出、收入、转账洞察
type InsightController struct {
	insightService *service.InsightService
}

// NewInsightController 创建数据洞察控制器实例
func NewInsightController(insightService *service.InsightService) *InsightController {
	return &InsightController{insightService: insightService}
}

// Expense 支出洞察
// GET /api/v1/insight/expense
// 分析支出数据，包括趋势、分类、账户等维度
// @Summary      Expense insight
// @Description  Analyze expense data including trends, categories, and accounts
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        start query string false "Start date (YYYY-MM-DD)"
// @Param        end query string false "End date (YYYY-MM-DD)"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/insight/expense [get]
// @Security     BearerAuth
func (c *InsightController) Expense(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	data, err := c.insightService.ExpenseInsight(userID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Income 收入洞察
// GET /api/v1/insight/income
// 分析收入数据，包括趋势、来源、账户等维度
// @Summary      Income insight
// @Description  Analyze income data including trends, sources, and accounts
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        start query string false "Start date (YYYY-MM-DD)"
// @Param        end query string false "End date (YYYY-MM-DD)"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/insight/income [get]
// @Security     BearerAuth
func (c *InsightController) Income(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	data, err := c.insightService.IncomeInsight(userID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Transfer 转账洞察
// GET /api/v1/insight/transfer
// 分析转账数据，包括账户间流动、频率等
// @Summary      Transfer insight
// @Description  Analyze transfer data including inter-account flows and frequency
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        start query string false "Start date (YYYY-MM-DD)"
// @Param        end query string false "End date (YYYY-MM-DD)"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/insight/transfer [get]
// @Security     BearerAuth
func (c *InsightController) Transfer(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	data, err := c.insightService.TransferInsight(userID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}
