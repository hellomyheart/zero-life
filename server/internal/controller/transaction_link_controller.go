// Package controller 提供HTTP请求处理控制器
// TransactionLinkController 交易关联控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type TransactionLinkController struct {
	service *service.TransactionLinkService
}

func NewTransactionLinkController(service *service.TransactionLinkService) *TransactionLinkController {
	return &TransactionLinkController{service: service}
}

func (ctrl *TransactionLinkController) Create(c *gin.Context) {
	var req request.CreateTransactionLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Create(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionLinkController) List(c *gin.Context) {
	var req request.TransactionLinkListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.List(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

func (ctrl *TransactionLinkController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}