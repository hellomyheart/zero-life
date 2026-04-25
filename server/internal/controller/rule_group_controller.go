// Package controller 提供HTTP请求处理控制器
// RuleGroupController 规则组控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// RuleGroupController 规则组管理控制器
// 处理规则组的创建、查询、更新、删除和批量执行等HTTP请求
// 规则组用于将多条规则组织在一起，支持按组批量执行
type RuleGroupController struct {
	ruleGroupService *service.RuleGroupService // 规则组业务服务
}

// NewRuleGroupController 创建规则组控制器实例
// 参数：
//   - ruleGroupService: 规则组业务服务实例
// 返回：
//   - *RuleGroupController: 规则组控制器实例
func NewRuleGroupController(ruleGroupService *service.RuleGroupService) *RuleGroupController {
	return &RuleGroupController{ruleGroupService: ruleGroupService}
}

// Create 创建规则组
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateRuleGroupReq（包含名称、排序、是否启用）
// 响应：创建成功的规则组信息
// @Summary      Create rule group
// @Description  Create a new rule group
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Param        body body request.CreateRuleGroupReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups [post]
// @Security     BearerAuth
func (ctrl *RuleGroupController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateRuleGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleGroupService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取规则组详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则组ID）
// 响应：规则组详细信息
// @Summary      Get rule group
// @Description  Get rule group details by ID
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule group ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups/{id} [get]
// @Security     BearerAuth
func (ctrl *RuleGroupController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.ruleGroupService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取规则组列表
// 返回当前用户的所有规则组
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：规则组列表
// @Summary      List rule groups
// @Description  Get all rule groups for current user
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups [get]
// @Security     BearerAuth
func (ctrl *RuleGroupController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.ruleGroupService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新规则组
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则组ID）
// 请求体：UpdateRuleGroupReq（包含名称、排序、是否启用）
// 响应：更新后的规则组信息
// @Summary      Update rule group
// @Description  Update rule group by ID
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule group ID"
// @Param        body body request.UpdateRuleGroupReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups/{id} [put]
// @Security     BearerAuth
func (ctrl *RuleGroupController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateRuleGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleGroupService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除规则组
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则组ID）
// 响应：删除成功返回nil
// 业务规则：规则组下有关联规则时无法删除
// @Summary      Delete rule group
// @Description  Delete rule group by ID (cannot delete if rules are attached)
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule group ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups/{id} [delete]
// @Security     BearerAuth
func (ctrl *RuleGroupController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.ruleGroupService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// ExecuteGroup 执行规则组
// 批量执行规则组下的所有规则，对指定日期范围内的交易应用规则
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（规则组ID）
// 请求体：ExecuteRuleGroupReq（包含start_date、end_date）
// 响应：RuleGroupExecuteResultResp（包含匹配数、成功数、失败数）
// @Summary      Execute rule group
// @Description  Execute all rules in a group for transactions in a date range
// @Tags         rule-groups
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Rule group ID"
// @Param        body body request.ExecuteRuleGroupReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/rule-groups/{id}/execute [post]
// @Security     BearerAuth
func (ctrl *RuleGroupController) ExecuteGroup(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.ExecuteRuleGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.ruleGroupService.ExecuteGroup(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
