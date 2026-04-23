package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type TransactionLinkController struct {
	txnLinkService *service.TransactionLinkService
}

func NewTransactionLinkController(txnLinkService *service.TransactionLinkService) *TransactionLinkController {
	return &TransactionLinkController{txnLinkService: txnLinkService}
}

func (ctrl *TransactionLinkController) Create(c *gin.Context) {
	var req request.CreateTransactionLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnLinkService.Create(&req)
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

	result, err := ctrl.txnLinkService.List(req.TransactionID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionLinkController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.txnLinkService.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
