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
func (ctrl *RecurringTransactionController) ProcessDue(c *gin.Context) {
	userID := c.GetUint64("user_id")

	created, err := ctrl.service.ProcessDue(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, gin.H{"created": created})
}