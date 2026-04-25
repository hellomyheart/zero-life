// Package controller 提供HTTP请求处理控制器
// CurrencyController 货币控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// CurrencyController 货币管理控制器
// 处理货币列表查询、启用/禁用、设置默认货币、汇率管理等HTTP请求
// 货币是系统级资源，不区分用户
type CurrencyController struct {
	currencyService *service.CurrencyService // 货币业务服务
}

// NewCurrencyController 创建货币控制器实例
// 参数：
//   - currencyService: 货币业务服务实例
// 返回：
//   - *CurrencyController: 货币控制器实例
func NewCurrencyController(currencyService *service.CurrencyService) *CurrencyController {
	return &CurrencyController{currencyService: currencyService}
}

// List 获取货币列表
// 返回系统中所有可用的货币信息
// 参数：
//   - c: Gin上下文
// 响应：货币列表
// @Summary      List currencies
// @Description  Get all available currencies
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/currencies [get]
// @Security     BearerAuth
func (ctrl *CurrencyController) List(c *gin.Context) {
	result, err := ctrl.currencyService.List()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// UpdateStatus 更新货币启用/禁用状态
// 参数：
//   - c: Gin上下文，包含URL路径参数id
// 路径参数：id（货币ID）
// 请求体：SetCurrencyStatusReq（包含is_enabled字段）
// 响应：更新成功返回nil
// 业务规则：正在被账户使用的货币不能禁用，默认货币不能禁用
// @Summary      Update currency status
// @Description  Enable or disable a currency
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Currency ID"
// @Param        body body request.SetCurrencyStatusReq true "set currency status request"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/currencies/{id}/status [put]
// @Security     BearerAuth
func (ctrl *CurrencyController) UpdateStatus(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.SetCurrencyStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.currencyService.UpdateStatus(id, req.IsEnabled); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// SetDefault 设置默认货币
// 将指定货币设为系统默认货币
// 参数：
//   - c: Gin上下文，包含URL路径参数id
// 路径参数：id（货币ID）
// 响应：设置成功返回nil
// @Summary      Set default currency
// @Description  Set a currency as the default currency
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Currency ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/currencies/{id}/default [put]
// @Security     BearerAuth
func (ctrl *CurrencyController) SetDefault(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.currencyService.SetDefault(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// GetExchangeRates 获取汇率列表
// 返回系统中所有货币对的汇率信息
// 参数：
//   - c: Gin上下文
// 响应：汇率列表
// @Summary      Get exchange rates
// @Description  Get all exchange rates
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/currencies/exchange-rates [get]
// @Security     BearerAuth
func (ctrl *CurrencyController) GetExchangeRates(c *gin.Context) {
	result, err := ctrl.currencyService.GetExchangeRates()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// SetExchangeRate 设置汇率
// 设置指定货币对之间的汇率
// 参数：
//   - c: Gin上下文
// 请求体：SetExchangeRateReq（包含源货币ID、目标货币ID、汇率值）
// 响应：设置成功返回nil
// @Summary      Set exchange rate
// @Description  Set exchange rate between two currencies
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        body body request.SetExchangeRateReq true "set exchange rate request"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/currencies/exchange-rates [post]
// @Security     BearerAuth
func (ctrl *CurrencyController) SetExchangeRate(c *gin.Context) {
	var req request.SetExchangeRateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.currencyService.SetExchangeRate(&req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
