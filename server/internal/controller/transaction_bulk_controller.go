package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type TransactionBulkController struct {
	bulkService *service.TransactionBulkService
}

func NewTransactionBulkController(bulkService *service.TransactionBulkService) *TransactionBulkController {
	return &TransactionBulkController{bulkService: bulkService}
}

func (ctrl *TransactionBulkController) BulkEdit(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.BulkEditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.bulkService.BulkEdit(userID, &req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *TransactionBulkController) BulkDelete(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.BulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.bulkService.BulkDelete(userID, &req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *TransactionBulkController) ConvertType(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.ConvertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.bulkService.ConvertType(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionBulkController) Clone(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.bulkService.CloneTransaction(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
