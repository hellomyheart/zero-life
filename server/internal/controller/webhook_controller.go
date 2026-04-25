// Package controller 提供HTTP请求处理控制器
// Webhook控制器，处理Webhook的增删改查
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// WebhookController Webhook管理控制器
// 处理Webhook的创建、查询、更新、删除和投递记录查询等HTTP请求
// Webhook用于在特定事件（如交易创建、账单支付等）发生时向外部URL发送通知
type WebhookController struct {
	service *service.WebhookService // Webhook业务服务
}

// NewWebhookController 创建Webhook控制器实例
// 参数：
//   - service: Webhook业务服务实例
// 返回：
//   - *WebhookController: Webhook控制器实例
func NewWebhookController(service *service.WebhookService) *WebhookController {
	return &WebhookController{service: service}
}

// Create 创建Webhook
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateWebhookReq（包含名称、URL、触发事件类型）
// 响应：创建成功的Webhook信息
// 业务规则：URL必须以https://开头，触发事件类型仅支持预定义值
// @Summary      Create webhook
// @Description  Create a new webhook for event notifications
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        body body request.CreateWebhookReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks [post]
// @Security     BearerAuth
func (ctrl *WebhookController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateWebhookReq
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

// Get 获取Webhook详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（Webhook ID）
// 响应：Webhook详细信息
// @Summary      Get webhook
// @Description  Get webhook details by ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Webhook ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks/{id} [get]
// @Security     BearerAuth
func (ctrl *WebhookController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取Webhook列表
// 返回当前用户的所有Webhook
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：Webhook列表
// @Summary      List webhooks
// @Description  Get all webhooks for current user
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks [get]
// @Security     BearerAuth
func (ctrl *WebhookController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.service.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新Webhook
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（Webhook ID）
// 请求体：UpdateWebhookReq（包含名称、URL、触发事件类型、是否启用）
// 响应：更新后的Webhook信息
// @Summary      Update webhook
// @Description  Update webhook by ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Webhook ID"
// @Param        body body request.UpdateWebhookReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks/{id} [put]
// @Security     BearerAuth
func (ctrl *WebhookController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateWebhookReq
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

// Delete 删除Webhook
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（Webhook ID）
// 响应：删除成功返回nil
// @Summary      Delete webhook
// @Description  Delete webhook by ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Webhook ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks/{id} [delete]
// @Security     BearerAuth
func (ctrl *WebhookController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// ListDeliveries 获取Webhook投递记录列表
// 查询指定Webhook的历史投递记录，包括HTTP状态码、错误信息和投递时间
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（Webhook ID）
// 查询参数：page（页码）、page_size（每页数量）
// 响应：投递记录列表
// @Summary      List webhook deliveries
// @Description  Get delivery history for a webhook including HTTP status, errors and timestamps
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Webhook ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/webhooks/{id}/deliveries [get]
// @Security     BearerAuth
func (ctrl *WebhookController) ListDeliveries(c *gin.Context) {
	userID := c.GetUint64("user_id")
	webhookID := parseIDParam(c, "id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := ctrl.service.ListDeliveries(userID, webhookID, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}