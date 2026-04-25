// Package controller 提供HTTP请求处理控制器
// RecurrenceController 周期性交易控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// RecurrenceController 周期性交易管理控制器
// 处理周期性交易的创建、查询、更新、删除和手动触发等HTTP请求
// 周期性交易用于自动创建重复的交易记录（如月薪、房租等）
type RecurrenceController struct {
	recurrenceService *service.RecurrenceService // 周期性交易业务服务
}

// NewRecurrenceController 创建周期性交易控制器实例
// 参数：
//   - recurrenceService: 周期性交易业务服务实例
// 返回：
//   - *RecurrenceController: 周期性交易控制器实例
func NewRecurrenceController(recurrenceService *service.RecurrenceService) *RecurrenceController {
	return &RecurrenceController{recurrenceService: recurrenceService}
}

// Create 创建周期性交易
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateRecurrenceReq（包含标题、类型、金额、重复频率、下次执行日期等）
// 响应：创建成功的周期性交易信息
func (ctrl *RecurrenceController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateRecurrenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.recurrenceService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取周期性交易详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（周期性交易ID）
// 响应：周期性交易详细信息
func (ctrl *RecurrenceController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.recurrenceService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取周期性交易列表（分页）
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：RecurrenceListReq（包含page、page_size）
// 响应：周期性交易分页列表
func (ctrl *RecurrenceController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.RecurrenceListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.recurrenceService.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新周期性交易
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（周期性交易ID）
// 请求体：UpdateRecurrenceReq（包含需要更新的字段）
// 响应：更新后的周期性交易信息
func (ctrl *RecurrenceController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateRecurrenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.recurrenceService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除周期性交易
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（周期性交易ID）
// 响应：删除成功返回nil
func (ctrl *RecurrenceController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.recurrenceService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// Trigger 手动触发周期性交易
// 手动执行一次周期性交易，创建对应的交易记录
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（周期性交易ID）
// 响应：创建的交易信息
func (ctrl *RecurrenceController) Trigger(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.recurrenceService.Trigger(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
