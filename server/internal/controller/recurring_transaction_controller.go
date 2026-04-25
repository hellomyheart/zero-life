// Package controller 提供HTTP请求处理控制器
// RecurringTransactionController 循环交易控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// RecurringTransactionController 循环交易管理控制器
// 处理循环交易的创建、查询、更新、删除和处理到期交易等HTTP请求
// 循环交易用于管理定期重复的交易（如月薪、房租等），到期时自动创建交易记录
type RecurringTransactionController struct {
	service *service.RecurringTransactionService // 循环交易业务服务
}

// NewRecurringTransactionController 创建循环交易控制器实例
// 参数：
//   - service: 循环交易业务服务实例
// 返回：
//   - *RecurringTransactionController: 循环交易控制器实例
func NewRecurringTransactionController(service *service.RecurringTransactionService) *RecurringTransactionController {
	return &RecurringTransactionController{service: service}
}

// Create 创建循环交易
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateRecurringTransactionReq（包含描述、金额、重复类型、开始日期等）
// 响应：创建成功的循环交易信息
// @Summary      Create recurring transaction
// @Description  Create a new recurring transaction
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Param        body body request.CreateRecurringTransactionReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions [post]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateRecurringTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取循环交易详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（循环交易ID）
// 响应：循环交易详细信息
// @Summary      Get recurring transaction
// @Description  Get recurring transaction details by ID
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Recurring transaction ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions/{id} [get]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取循环交易列表（分页）
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：RecurringTransactionListReq（包含page、page_size、active过滤）
// 响应：循环交易分页列表
// @Summary      List recurring transactions
// @Description  Get paginated list of recurring transactions
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        active query bool false "Filter by active status"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions [get]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.RecurringTransactionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新循环交易
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（循环交易ID）
// 请求体：UpdateRecurringTransactionReq（包含需要更新的字段）
// 响应：更新后的循环交易信息
// @Summary      Update recurring transaction
// @Description  Update recurring transaction by ID
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Recurring transaction ID"
// @Param        body body request.UpdateRecurringTransactionReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions/{id} [put]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateRecurringTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除循环交易
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（循环交易ID）
// 响应：删除成功返回nil
// @Summary      Delete recurring transaction
// @Description  Delete recurring transaction by ID
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Recurring transaction ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions/{id} [delete]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// ProcessDue 处理到期循环交易
// 批量检查并处理所有到期的循环交易，自动创建对应的交易记录
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：创建的交易数量（{"created": N}）
// @Summary      Process due recurring transactions
// @Description  Check and process all due recurring transactions, creating corresponding transaction records
// @Tags         recurring-transactions
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/recurring-transactions/process-due [post]
// @Security     BearerAuth
func (ctrl *RecurringTransactionController) ProcessDue(c *gin.Context) {
	userID := c.GetUint64("user_id")

	created, err := ctrl.service.ProcessDue(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, gin.H{"created": created})
}