// Package controller 提供HTTP请求处理控制器
// LinkTypeController 关联类型控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type LinkTypeController struct {
	linkTypeService *service.LinkTypeService
}

func NewLinkTypeController(linkTypeService *service.LinkTypeService) *LinkTypeController {
	return &LinkTypeController{linkTypeService: linkTypeService}
}

func (ctrl *LinkTypeController) Create(c *gin.Context) {
	var req request.CreateLinkTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.Create(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *LinkTypeController) Get(c *gin.Context) {
	id := parseIDParam(c, "id")

	result, err := ctrl.linkTypeService.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *LinkTypeController) List(c *gin.Context) {
	var req request.LinkTypeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.List(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

func (ctrl *LinkTypeController) Update(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.UpdateLinkTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.Update(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *LinkTypeController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.linkTypeService.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
