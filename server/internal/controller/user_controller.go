// Package controller 控制器层
// 处理HTTP请求，调用Service层业务逻辑
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// UserController 用户控制器
// 处理用户管理相关请求（管理员功能）
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// List 获取用户列表
// GET /api/v1/users
// 仅管理员可访问
// @Summary      List users
// @Description  Get paginated user list (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users [get]
// @Security     BearerAuth
func (c *UserController) List(ctx *gin.Context) {
	var req request.UserListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := c.userService.List(&req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	SuccessPage(ctx, result)
}

// Get 获取用户详情
// GET /api/v1/users/:id
// 仅管理员可访问
// @Summary      Get user
// @Description  Get user details by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id} [get]
// @Security     BearerAuth
func (c *UserController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := c.userService.Get(id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, result)
}

// Update 更新用户信息
// PUT /api/v1/users/:id
// 仅管理员可访问
// @Summary      Update user
// @Description  Update user info by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Param        body body request.UpdateUserReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id} [put]
// @Security     BearerAuth
func (c *UserController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	var req request.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := c.userService.Update(id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, result)
}

// Delete 删除用户
// DELETE /api/v1/users/:id
// 仅管理员可访问
// @Summary      Delete user
// @Description  Delete user by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id} [delete]
// @Security     BearerAuth
func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.userService.Delete(id); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// ChangeRole 修改用户角色
// PUT /api/v1/users/:id/role
// 仅管理员可访问
// @Summary      Change user role
// @Description  Change user role by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Param        body body request.ChangeRoleReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id}/role [put]
// @Security     BearerAuth
func (c *UserController) ChangeRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	var req request.ChangeRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.userService.ChangeRole(id, req.Role); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// Lock 锁定用户
// POST /api/v1/users/:id/lock
// 仅管理员可访问
// @Summary      Lock user
// @Description  Lock user account by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id}/lock [post]
// @Security     BearerAuth
func (c *UserController) Lock(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.userService.Lock(id); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// Unlock 解锁用户
// POST /api/v1/users/:id/unlock
// 仅管理员可访问
// @Summary      Unlock user
// @Description  Unlock user account by ID (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/users/{id}/unlock [post]
// @Security     BearerAuth
func (c *UserController) Unlock(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.userService.Unlock(id); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}
