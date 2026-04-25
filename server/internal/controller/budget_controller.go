// Package controller 提供HTTP请求处理控制器
// BudgetController 预算控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// BudgetController 预算管理控制器
// 处理预算的创建、查询、更新、删除和历史记录等HTTP请求
// 预算用于控制指定分类的支出金额，支持月度和年度周期
type BudgetController struct {
	budgetService *service.BudgetService // 预算业务服务
}

// NewBudgetController 创建预算控制器实例
// 参数：
//   - budgetService: 预算业务服务实例
// 返回：
//   - *BudgetController: 预算控制器实例
func NewBudgetController(budgetService *service.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

// Create 创建预算
// 接收创建预算请求，设置预算名称、金额、周期和关联分类
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateBudgetReq（包含名称、金额、周期、分类ID列表）
// 响应：创建成功的预算信息
// 业务规则：预算金额必须大于0，至少关联一个分类
// @Summary      Create budget
// @Description  Create a new budget for the current user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Param        body body request.CreateBudgetReq true "create budget request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets [post]
// @Security     BearerAuth
func (ctrl *BudgetController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateBudgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.budgetService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取预算详情
// 根据URL路径中的预算ID，查询预算详细信息（含已支出金额、剩余金额、使用率）
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（预算ID）
// 响应：预算详细信息
// @Summary      Get budget
// @Description  Get budget details by ID
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Budget ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets/{id} [get]
// @Security     BearerAuth
func (ctrl *BudgetController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.budgetService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取预算列表
// 返回当前用户的所有预算，含已支出金额、剩余金额和使用率
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：预算列表
// @Summary      List budgets
// @Description  Get all budgets for the current user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets [get]
// @Security     BearerAuth
func (ctrl *BudgetController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.budgetService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新预算
// 根据URL路径中的预算ID和请求体中的更新字段，修改预算信息
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（预算ID）
// 请求体：UpdateBudgetReq（包含需要更新的字段）
// 响应：更新后的预算信息
// @Summary      Update budget
// @Description  Update budget by ID
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Budget ID"
// @Param        body body request.UpdateBudgetReq true "update budget request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets/{id} [put]
// @Security     BearerAuth
func (ctrl *BudgetController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateBudgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.budgetService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除预算
// 根据URL路径中的预算ID删除预算
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（预算ID）
// 响应：删除成功返回nil
// @Summary      Delete budget
// @Description  Delete budget by ID
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Budget ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets/{id} [delete]
// @Security     BearerAuth
func (ctrl *BudgetController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.budgetService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// GetHistory 获取预算历史记录
// 查询指定预算的历史执行情况，包括每个周期的预算金额、实际支出和使用率
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（预算ID）
// 响应：预算历史记录列表
// @Summary      Get budget history
// @Description  Get budget history by ID
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Budget ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/budgets/{id}/history [get]
// @Security     BearerAuth
func (ctrl *BudgetController) GetHistory(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.budgetService.GetHistory(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
