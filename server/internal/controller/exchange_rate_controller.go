// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ExchangeRateController 汇率控制器
// 处理汇率相关的HTTP请求，包括汇率的创建、查询、更新、删除和货币转换
type ExchangeRateController struct {
	rateService *service.ExchangeRateService
}

// NewExchangeRateController 创建汇率控制器实例
// 参数：
//   - rateService: 汇率业务服务
// 返回：
//   - *ExchangeRateController: 汇率控制器实例
func NewExchangeRateController(rateService *service.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{rateService: rateService}
}

// Create 创建汇率
// 接收创建汇率请求，调用业务层创建汇率记录
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateExchangeRateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.rateService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取汇率详情
// 根据ID获取汇率信息
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.rateService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取汇率列表
// 支持按货币对、日期范围过滤和分页
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ExchangeRateListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.rateService.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新汇率
// 更新指定ID的汇率信息
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateExchangeRateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.rateService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除汇率
// 删除指定ID的汇率记录
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.rateService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// GetLatest 获取最新汇率
// 获取指定货币对的最新汇率
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) GetLatest(c *gin.Context) {
	userID := c.GetUint64("user_id")

	fromCurrencyID := parseIDParam(c, "from")
	toCurrencyID := parseIDParam(c, "to")

	result, err := ctrl.rateService.GetLatest(userID, fromCurrencyID, toCurrencyID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Convert 货币转换
// 根据汇率将金额从源货币转换为目标货币
// 参数：
//   - c: Gin上下文
func (ctrl *ExchangeRateController) Convert(c *gin.Context) {
	userID := c.GetUint64("user_id")

	fromCurrencyID := parseIDParam(c, "from")
	toCurrencyID := parseIDParam(c, "to")
	amount := c.Query("amount")

	if amount == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	var date *time.Time
	if dateStr := c.Query("date"); dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
			return
		}
		date = &t
	}

	result, err := ctrl.rateService.Convert(userID, fromCurrencyID, toCurrencyID, amount, date)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, gin.H{"amount": result})
}
