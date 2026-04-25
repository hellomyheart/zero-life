// Package controller 提供HTTP请求处理控制器
// AdminController 管理员控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// AdminController 管理员配置管理控制器
// 处理系统配置的查询、更新和邮件测试等HTTP请求
// 仅管理员可访问，用于管理系统级配置参数
type AdminController struct {
	adminService *service.AdminService // 管理员业务服务
}

// NewAdminController 创建管理员控制器实例
// 参数：
//   - adminService: 管理员业务服务实例
// 返回：
//   - *AdminController: 管理员控制器实例
func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

// GetConfiguration 获取指定配置项
// 根据配置名称查询单个系统配置的值
// 参数：
//   - c: Gin上下文
// 路径参数：name（配置名称）
// 响应：配置详细信息
// 业务规则：name不能为空
// @Summary      Get configuration by name
// @Description  Get a single system configuration value by name
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        name path string true "Configuration name"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/configurations/{name} [get]
// @Security     BearerAuth
func (ctrl *AdminController) GetConfiguration(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.GetConfiguration(name)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// ListConfigurations 获取所有配置项列表
// 返回系统中所有配置项
// 参数：
//   - c: Gin上下文
// 响应：配置项列表
// @Summary      List all configurations
// @Description  Get all system configuration items
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/configurations [get]
// @Security     BearerAuth
func (ctrl *AdminController) ListConfigurations(c *gin.Context) {
	result, err := ctrl.adminService.ListConfigurations()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// UpdateConfiguration 更新指定配置项
// 根据配置名称更新配置值
// 参数：
//   - c: Gin上下文
// 路径参数：name（配置名称）
// 请求体：AdminUpdateConfigurationReq（包含value字段）
// 响应：更新后的配置信息
// @Summary      Update configuration
// @Description  Update a system configuration value by name
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        name path string true "Configuration name"
// @Param        body body request.AdminUpdateConfigurationReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/configurations/{name} [put]
// @Security     BearerAuth
func (ctrl *AdminController) UpdateConfiguration(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	var req request.AdminUpdateConfigurationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.UpdateConfiguration(name, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// TestEmail 测试邮件发送
// 向指定邮箱发送测试邮件，验证SMTP配置是否正确
// 参数：
//   - c: Gin上下文
// 请求体：AdminTestEmailReq（包含email字段）
// 响应：AdminTestEmailResp（包含success和message）
// @Summary      Test email sending
// @Description  Send a test email to verify SMTP configuration
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body request.AdminTestEmailReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/admin/test-email [post]
// @Security     BearerAuth
func (ctrl *AdminController) TestEmail(c *gin.Context) {
	var req request.AdminTestEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.TestEmail(req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
