// Package controller 提供HTTP请求处理控制器
// RuleController 规则控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// RuleController 规则管理控制器
// 处理规则的创建、查询、更新、删除、启用/禁用切换和执行等HTTP请求
// 规则用于在交易创建或更新时自动执行操作（如自动分类、添加标签等）
type RuleController struct {
	ruleService *service.RuleService // 规则业务服务
}

// NewRuleController 创建规则控制器实例
// 参数：
//   - ruleService: 规则业务服务实例
// 返回：
//   - *RuleController: 规则控制器实例
func NewRuleController(ruleService *service.RuleService) *RuleController {
	return &RuleController{ruleService: ruleService}
}

// Create 创建规则
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateRuleReq（包含名称、优先级、逻辑类型、触发条件、条件列表、动作列表）
// 响应：创建成功的规则信息
// 业务规则：至少需要一个条件和一个动作，逻辑类型支持and/or
// @Summary      Create rule
// @Description  Create a new rule for the current user
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        body body request.CreateRuleReq true "create rule request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules [post]
// @Security     BearerAuth
func (ctrl *RuleController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取规则详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则ID）
// 响应：规则详细信息（含条件和动作列表）
// @Summary      Get rule
// @Description  Get rule details by ID
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules/{id} [get]
// @Security     BearerAuth
func (ctrl *RuleController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.ruleService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取规则列表
// 返回当前用户的所有规则
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：规则列表
// @Summary      List rules
// @Description  Get all rules for the current user
// @Tags         rules
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules [get]
// @Security     BearerAuth
func (ctrl *RuleController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.ruleService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新规则
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则ID）
// 请求体：UpdateRuleReq（包含需要更新的字段）
// 响应：更新后的规则信息
// @Summary      Update rule
// @Description  Update rule by ID
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule ID"
// @Param        body body request.UpdateRuleReq true "update rule request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules/{id} [put]
// @Security     BearerAuth
func (ctrl *RuleController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除规则
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则ID）
// 响应：删除成功返回nil
// @Summary      Delete rule
// @Description  Delete rule by ID
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules/{id} [delete]
// @Security     BearerAuth
func (ctrl *RuleController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.ruleService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// ToggleStatus 切换规则启用/禁用状态
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则ID）
// 响应：更新后的规则信息
// @Summary      Toggle rule status
// @Description  Toggle rule enabled/disabled status
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules/{id}/toggle [put]
// @Security     BearerAuth
func (ctrl *RuleController) ToggleStatus(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.ruleService.ToggleStatus(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Execute 执行规则
// 手动触发规则执行，对指定日期范围内的交易应用规则
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则ID）
// 请求体：ExecuteRuleReq（包含start_date、end_date）
// 响应：RuleExecuteResultResp（包含匹配数、成功数、失败数）
// @Summary      Execute rule
// @Description  Manually execute a rule on transactions within a date range
// @Tags         rules
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule ID"
// @Param        body body request.ExecuteRuleReq true "execute rule request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rules/{id}/execute [post]
// @Security     BearerAuth
func (ctrl *RuleController) Execute(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.ExecuteRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleService.Execute(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
