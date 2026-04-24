// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// AccountController 账户管理控制器，处理账户相关的HTTP请求
type AccountController struct {
	accountService *service.AccountService
}

// NewAccountController 创建账户控制器实例
func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

// Create 创建账户
func (ctrl *AccountController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.accountService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取账户详情
func (ctrl *AccountController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.accountService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取账户列表（分页）
func (ctrl *AccountController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.AccountListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.accountService.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新账户
func (ctrl *AccountController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.accountService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除账户
func (ctrl *AccountController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.accountService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// parseIDParam 从URL路径参数中解析ID
func parseIDParam(c *gin.Context, param string) uint64 {
	str := c.Param(param)
	if str == "" {
		return 0
	}
	if str[0] == '/' {
		str = str[1:]
	}
	id, _ := strconv.ParseUint(str, 10, 64)
	return id
}
