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

type RuleGroupController struct {
	ruleGroupService *service.RuleGroupService
}

func NewRuleGroupController(ruleGroupService *service.RuleGroupService) *RuleGroupController {
	return &RuleGroupController{ruleGroupService: ruleGroupService}
}

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

func (ctrl *RuleGroupController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.ruleGroupService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

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

func (ctrl *RuleGroupController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.ruleGroupService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

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
