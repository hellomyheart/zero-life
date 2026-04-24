// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ProfileController 用户资料控制器
// 处理用户资料相关的HTTP请求
type ProfileController struct {
	profileService *service.ProfileService
}

// NewProfileController 创建用户资料控制器实例
// 参数：
//   - profileService: 用户资料业务服务
// 返回：
//   - *ProfileController: 用户资料控制器实例
func NewProfileController(profileService *service.ProfileService) *ProfileController {
	return &ProfileController{profileService: profileService}
}

// GetProfile 获取用户资料
// 参数：
//   - c: Gin上下文
func (ctrl *ProfileController) GetProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.profileService.GetProfile(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// UpdateProfile 更新用户资料
// 参数：
//   - c: Gin上下文
func (ctrl *ProfileController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.profileService.UpdateProfile(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// ChangePassword 修改密码
// 参数：
//   - c: Gin上下文
func (ctrl *ProfileController) ChangePassword(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.profileService.ChangePassword(userID, &req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
