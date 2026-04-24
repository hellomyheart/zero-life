package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type TransactionController struct {
	txnService *service.TransactionService
}

func NewTransactionController(txnService *service.TransactionService) *TransactionController {
	return &TransactionController{txnService: txnService}
}

func (ctrl *TransactionController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.txnService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.TransactionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnService.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

func (ctrl *TransactionController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *TransactionController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.txnService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *TransactionController) Search(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.TransactionSearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnService.Search(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Split 拆分交易
// POST /api/v1/transactions/:id/split
// 将一笔交易拆分为多笔子交易
func (ctrl *TransactionController) Split(c *gin.Context) {
	userID := c.GetUint64("user_id")
	parentID := parseIDParam(c, "id")

	var req request.SplitTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.txnService.Split(userID, parentID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// GetSplits 获取拆分交易列表
// GET /api/v1/transactions/:id/splits
func (ctrl *TransactionController) GetSplits(c *gin.Context) {
	userID := c.GetUint64("user_id")
	parentID := parseIDParam(c, "id")

	result, err := ctrl.txnService.GetSplits(userID, parentID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// MergeSplits 合并拆分交易
// POST /api/v1/transactions/:id/merge
// 将拆分的子交易合并回父交易
func (ctrl *TransactionController) MergeSplits(c *gin.Context) {
	userID := c.GetUint64("user_id")
	parentID := parseIDParam(c, "id")

	if err := ctrl.txnService.MergeSplits(userID, parentID); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
