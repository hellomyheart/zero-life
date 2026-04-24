// Package controller 提供HTTP请求处理控制器
// AdminUserController 管理员用户管理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type AdminUserController struct {
	adminService *service.AdminService
}

func NewAdminUserController(adminService *service.AdminService) *AdminUserController {
	return &AdminUserController{adminService: adminService}
}

func (ctrl *AdminUserController) ListUsers(c *gin.Context) {
	var req request.AdminListUsersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	items, total, err := ctrl.adminService.ListUsers(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()
	SuccessPage(c, pagination.NewResult(items, total, params))
}

func (ctrl *AdminUserController) UpdateUser(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.AdminUpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.UpdateUser(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AdminUserController) DeleteUser(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.adminService.DeleteUser(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *AdminUserController) InviteUser(c *gin.Context) {
	var req request.AdminInviteUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.InviteUser(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
