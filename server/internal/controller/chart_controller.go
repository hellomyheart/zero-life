// Package controller 控制器层
// 处理HTTP请求，调用Service层业务逻辑
package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ChartController 图表控制器
// 提供各类图表数据接口，用于前端可视化展示
type ChartController struct {
	chartService *service.ChartService
}

// NewChartController 创建图表控制器实例
func NewChartController(chartService *service.ChartService) *ChartController {
	return &ChartController{chartService: chartService}
}

// Account 账户余额图表数据
// GET /api/v1/chart/account/:id
// 返回指定账户的余额变化趋势数据
func (c *ChartController) Account(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	accountIDStr := ctx.Param("id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // 默认最近一个月
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	data, err := c.chartService.AccountBalance(userID, accountID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Budget 预算图表数据
// GET /api/v1/chart/budget/:id
// 返回指定预算的支出趋势数据
func (c *ChartController) Budget(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	budgetIDStr := ctx.Param("id")
	budgetID, err := strconv.ParseUint(budgetIDStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, err := c.chartService.BudgetSpending(userID, budgetID)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Category 分类图表数据
// GET /api/v1/chart/category
// 返回分类支出/收入分布数据
func (c *ChartController) Category(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")
	chartType := ctx.Query("type") // expense 或 income

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	data, err := c.chartService.CategoryDistribution(userID, startDate, endDate, chartType)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Tag 标签图表数据
// GET /api/v1/chart/tag
// 返回标签支出/收入分布数据
func (c *ChartController) Tag(ctx *gin.Context) {
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

	data, err := c.chartService.TagDistribution(userID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}

// Transaction 交易图表数据
// GET /api/v1/chart/transaction
// 返回交易趋势数据（收入、支出、转账）
func (c *ChartController) Transaction(ctx *gin.Context) {
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

	data, err := c.chartService.TransactionTrend(userID, startDate, endDate)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	Success(ctx, data)
}
