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

// AdminUserController 管理员用户管理控制器
// 处理用户列表查询、用户信息更新、用户删除和邀请用户等HTTP请求
// 仅管理员可访问，用于管理系统中所有用户
type AdminUserController struct {
	adminService *service.AdminService // 管理员业务服务
}

// NewAdminUserController 创建管理员用户管理控制器实例
// 参数：
//   - adminService: 管理员业务服务实例
// 返回：
//   - *AdminUserController: 管理员用户管理控制器实例
func NewAdminUserController(adminService *service.AdminService) *AdminUserController {
	return &AdminUserController{adminService: adminService}
}

// ListUsers 获取用户列表（分页）
// 支持按关键词搜索用户
// 参数：
//   - c: Gin上下文
// 查询参数：AdminListUsersReq（包含page、page_size、search）
// 响应：用户分页列表
// @Summary      List users
// @Description  Get paginated user list with optional search
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        search query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/users [get]
// @Security     BearerAuth
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

// UpdateUser 更新用户信息
// 管理员更新指定用户的信息（昵称、角色、语言、时区）
// 参数：
//   - c: Gin上下文
// 路径参数：id（用户ID）
// 请求体：AdminUpdateUserReq（包含昵称、角色、语言、时区）
// 响应：更新后的用户信息
// @Summary      Update user
// @Description  Admin update user info (nickname, role, language, timezone)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Param        body body request.AdminUpdateUserReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/users/{id} [put]
// @Security     BearerAuth
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

// DeleteUser 删除用户
// 管理员删除指定用户
// 参数：
//   - c: Gin上下文
// 路径参数：id（用户ID）
// 响应：删除成功返回nil
// @Summary      Delete user
// @Description  Admin delete a user
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/users/{id} [delete]
// @Security     BearerAuth
func (ctrl *AdminUserController) DeleteUser(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.adminService.DeleteUser(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// InviteUser 邀请用户
// 管理员通过邮箱邀请新用户注册，发送邀请邮件
// 参数：
//   - c: Gin上下文
// 请求体：AdminInviteUserReq（包含邮箱、昵称、角色）
// 响应：邀请结果
// @Summary      Invite user
// @Description  Admin invite a new user by email
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body request.AdminInviteUserReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/users/invite [post]
// @Security     BearerAuth
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
