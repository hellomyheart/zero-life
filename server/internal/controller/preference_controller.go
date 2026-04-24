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

type PreferenceController struct {
	service *service.PreferenceService
}

func NewPreferenceController(service *service.PreferenceService) *PreferenceController {
	return &PreferenceController{service: service}
}

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

func (ctrl *PreferenceController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.service.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *PreferenceController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	key := c.Param("key")

	if err := ctrl.service.Delete(userID, key); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}