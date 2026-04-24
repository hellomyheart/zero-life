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
func (c *UserController) List(ctx *gin.Context) {
	// 检查是否为管理员
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) Get(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) Update(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) Delete(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) ChangeRole(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) Lock(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
func (c *UserController) Unlock(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		Error(ctx, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

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
