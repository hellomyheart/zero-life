// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ReconciliationController 对账控制器
// 负责处理对账相关的HTTP请求，包括对账记录的创建、查询、更新和删除
// 所有操作均需传入用户ID，确保用户只能访问自己的对账数据
type ReconciliationController struct {
	service *service.ReconciliationService
}

// NewReconciliationController 创建对账控制器实例
// 参数：
//   - service: 对账业务服务实例
// 返回：
//   - *ReconciliationController: 对账控制器实例
func NewReconciliationController(service *service.ReconciliationService) *ReconciliationController {
	return &ReconciliationController{service: service}
}

// Create 创建对账记录
// 从上下文中提取用户ID，绑定请求参数，调用服务层创建对账记录
// 请求方法：POST
// 请求路径：/reconciliations
// 请求体：CreateReconciliationReq（JSON格式）
func (ctrl *ReconciliationController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateReconciliationReq
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

// Get 获取对账记录详情
// 从上下文中提取用户ID，确保只能查询属于当前用户的对账记录
// 请求方法：GET
// 请求路径：/reconciliations/:id
func (ctrl *ReconciliationController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取对账记录列表（分页）
// 从上下文中提取用户ID，确保只能查询属于当前用户的对账记录
// 请求方法：GET
// 请求路径：/reconciliations
// 查询参数：ReconciliationListReq
func (ctrl *ReconciliationController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReconciliationListReq
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

// Update 更新对账记录
// 从上下文中提取用户ID，确保只能更新属于当前用户的对账记录
// 请求方法：PUT
// 请求路径：/reconciliations/:id
// 请求体：UpdateReconciliationReq（JSON格式）
func (ctrl *ReconciliationController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateReconciliationReq
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

// Delete 删除对账记录
// 从上下文中提取用户ID，确保只能删除属于当前用户的对账记录
// 请求方法：DELETE
// 请求路径：/reconciliations/:id
func (ctrl *ReconciliationController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
