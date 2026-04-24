// Package controller 提供HTTP请求处理控制器
// TagController 标签控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type TagController struct {
	tagService *service.TagService
}

func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (ctrl *TagController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.tagService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TagController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.tagService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TagController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.tagService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TagController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.tagService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
