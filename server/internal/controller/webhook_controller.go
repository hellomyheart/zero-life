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

type WebhookController struct {
	service *service.WebhookService
}

func NewWebhookController(service *service.WebhookService) *WebhookController {
	return &WebhookController{service: service}
}

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

func (ctrl *WebhookController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.service.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

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

func (ctrl *WebhookController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

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