// Package controller 提供HTTP请求处理控制器
// ObjectGroupController 对象分组控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ObjectGroupController 对象分组管理控制器
// 处理对象分组的创建、查询、更新和删除等HTTP请求
// 对象分组用于将交易、账单、预算、存钱罐等对象进行分组管理
type ObjectGroupController struct {
	service *service.ObjectGroupService // 对象分组业务服务
}

// NewObjectGroupController 创建对象分组控制器实例
// 参数：
//   - service: 对象分组业务服务实例
// 返回：
//   - *ObjectGroupController: 对象分组控制器实例
func NewObjectGroupController(service *service.ObjectGroupService) *ObjectGroupController {
	return &ObjectGroupController{service: service}
}

// Create 创建对象分组
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateObjectGroupReq（包含名称、对象类型、对象ID）
// 响应：创建成功的分组信息
func (ctrl *ObjectGroupController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateObjectGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取对象分组详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（分组ID）
// 响应：分组详细信息
func (ctrl *ObjectGroupController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取对象分组列表
// 根据对象类型过滤分组列表
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：groupable_type（对象类型，如transaction/bill/budget/piggy_bank）
// 响应：分组列表
func (ctrl *ObjectGroupController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")
	groupableType := c.Query("groupable_type")

	result, err := ctrl.service.List(userID, groupableType)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新对象分组
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（分组ID）
// 请求体：UpdateObjectGroupReq（包含名称）
// 响应：更新后的分组信息
func (ctrl *ObjectGroupController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateObjectGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除对象分组
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（分组ID）
// 响应：删除成功返回nil
func (ctrl *ObjectGroupController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}