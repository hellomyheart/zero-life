// Package controller 提供HTTP请求处理控制器
// TransactionController 交易控制器，处理交易的增删改查和拆分合并
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// TransactionController 交易控制器
// 处理交易相关的HTTP请求，包括交易的创建、查询、更新、删除、拆分和合并
type TransactionController struct {
	txnService *service.TransactionService // 交易业务服务
}

// NewTransactionController 创建交易控制器实例
// 参数：
//   - txnService: 交易业务服务
// 返回：
//   - *TransactionController: 交易控制器实例
func NewTransactionController(txnService *service.TransactionService) *TransactionController {
	return &TransactionController{txnService: txnService}
}

// Create 创建交易
// 接收创建交易请求，调用业务层创建交易记录
// 参数：
//   - c: Gin上下文，包含请求信息
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
