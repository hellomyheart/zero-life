// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// AuthController 认证管理控制器，处理注册、登录、Token刷新、密码重置等HTTP请求
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register 用户注册
func (ctrl *AuthController) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.Register(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.Login(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// RefreshToken 刷新访问令牌
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// ForgotPassword 忘记密码，发送重置邮件
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ForgotPassword(req.Email); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// ResetPassword 重置密码
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ResetPassword(req.Token, req.Password); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// GetProfile 获取当前用户个人信息
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// UpdateProfile 更新当前用户个人信息
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.UpdateProfile(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// ChangePassword 修改密码
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ChangePassword(userID, &req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// handleError 统一错误处理，根据错误码映射HTTP状态码
func handleError(c *gin.Context, err error) {
	if e, ok := err.(*errcode.Error); ok {
		status := http.StatusInternalServerError
		switch e.Code {
		case 400:
			status = http.StatusBadRequest
		case 401:
			status = http.StatusUnauthorized
		case 403:
			status = http.StatusForbidden
		case 404:
			status = http.StatusNotFound
		case 429:
			status = http.StatusTooManyRequests
		default:
			if e.Code >= 10000 {
				status = http.StatusBadRequest
			}
		}
		Error(c, status, e)
		return
	}
	Error(c, http.StatusInternalServerError, errcode.ErrInternal)
}
