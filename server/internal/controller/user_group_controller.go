// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// UserGroupController 用户组控制器
type UserGroupController struct {
	groupService *service.UserGroupService
}

// NewUserGroupController 创建用户组控制器实例
func NewUserGroupController(groupService *service.UserGroupService) *UserGroupController {
	return &UserGroupController{groupService: groupService}
}

// Create 创建用户组
func (ctrl *UserGroupController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateUserGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.groupService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取用户组详情
func (ctrl *UserGroupController) Get(c *gin.Context) {
	id := parseIDParam(c, "id")

	result, err := ctrl.groupService.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取用户组列表
func (ctrl *UserGroupController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.UserGroupListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.groupService.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新用户组
func (ctrl *UserGroupController) Update(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.UpdateUserGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.groupService.Update(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除用户组
func (ctrl *UserGroupController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.groupService.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
