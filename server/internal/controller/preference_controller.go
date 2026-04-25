// Package controller 提供HTTP请求处理控制器
// PreferenceController 偏好设置控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// PreferenceController 偏好设置管理控制器
// 处理用户偏好设置的获取、设置、列表查询和删除等HTTP请求
// 偏好设置用于存储用户的个性化配置（如默认货币、首页布局等）
type PreferenceController struct {
	service *service.PreferenceService // 偏好设置业务服务
}

// NewPreferenceController 创建偏好设置控制器实例
// 参数：
//   - service: 偏好设置业务服务实例
// 返回：
//   - *PreferenceController: 偏好设置控制器实例
func NewPreferenceController(service *service.PreferenceService) *PreferenceController {
	return &PreferenceController{service: service}
}

// Get 获取偏好设置
// 根据key查询当前用户的偏好设置值
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：key（偏好设置键名）
// 响应：偏好设置值
// 业务规则：key不能为空
func (ctrl *PreferenceController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	key := c.Query("key")

	if key == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Get(userID, key)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Set 设置偏好设置
// 设置当前用户的偏好键值对，已存在则更新
// 参数：
//   - c: Gin上下文，包含用户身份
// 请求体：SetPreferenceReq（包含key和value）
// 响应：设置后的偏好信息
func (ctrl *PreferenceController) Set(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.SetPreferenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Set(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取所有偏好设置
// 返回当前用户的所有偏好设置列表
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：偏好设置列表
func (ctrl *PreferenceController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.service.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除偏好设置
// 根据key删除当前用户的偏好设置
// 参数：
//   - c: Gin上下文，包含用户身份
// 路径参数：key（偏好设置键名）
// 响应：删除成功返回nil
func (ctrl *PreferenceController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	key := c.Param("key")

	if err := ctrl.service.Delete(userID, key); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}