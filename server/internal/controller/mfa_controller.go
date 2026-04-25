// Package controller 控制器层
// 处理HTTP请求，调用Service层业务逻辑
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// MFAController MFA控制器
// 处理多因素认证相关请求
type MFAController struct {
	mfaService *service.MFAService
}

// NewMFAController 创建MFA控制器实例
func NewMFAController(mfaService *service.MFAService) *MFAController {
	return &MFAController{mfaService: mfaService}
}

// Setup 初始化MFA设置
// POST /api/v1/mfa/setup
// 生成MFA密钥和二维码URL
// @Summary      Setup MFA
// @Description  Generate MFA secret key and QR code URL
// @Tags         mfa
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/mfa/setup [post]
// @Security     BearerAuth
func (c *MFAController) Setup(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	result, err := c.mfaService.Setup(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, result)
}

// Enable 启用MFA
// POST /api/v1/mfa/enable
// 验证MFA代码并启用多因素认证
// @Summary      Enable MFA
// @Description  Verify MFA code and enable multi-factor authentication
// @Tags         mfa
// @Accept       json
// @Produce      json
// @Param        body body request.MFAVerifyReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/mfa/enable [post]
// @Security     BearerAuth
func (c *MFAController) Enable(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.MFAVerifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.mfaService.Enable(userID, req.Code); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// Disable 禁用MFA
// POST /api/v1/mfa/disable
// 验证MFA代码并禁用多因素认证
// @Summary      Disable MFA
// @Description  Verify MFA code and disable multi-factor authentication
// @Tags         mfa
// @Accept       json
// @Produce      json
// @Param        body body request.MFAVerifyReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/mfa/disable [post]
// @Security     BearerAuth
func (c *MFAController) Disable(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.MFAVerifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.mfaService.Disable(userID, req.Code); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// Verify 验证MFA代码
// POST /api/v1/mfa/verify
// 用于登录后的MFA验证
// @Summary      Verify MFA code
// @Description  Verify MFA code for post-login authentication
// @Tags         mfa
// @Accept       json
// @Produce      json
// @Param        body body request.MFAVerifyReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/mfa/verify [post]
// @Security     BearerAuth
func (c *MFAController) Verify(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.MFAVerifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := c.mfaService.Verify(userID, req.Code)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, result)
}

// Status 获取MFA状态
// GET /api/v1/mfa/status
// @Summary      Get MFA status
// @Description  Get current user's MFA enabled status
// @Tags         mfa
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/mfa/status [get]
// @Security     BearerAuth
func (c *MFAController) Status(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	result, err := c.mfaService.Status(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, result)
}
