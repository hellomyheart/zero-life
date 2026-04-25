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
// @Summary      Create transaction
// @Description  Create a new transaction for the current user
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body body request.CreateTransactionReq true "create transaction request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions [post]
// @Security     BearerAuth
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

// @Summary      Get transaction
// @Description  Get transaction details by ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id} [get]
// @Security     BearerAuth
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

// @Summary      List transactions
// @Description  Get paginated list of transactions for the current user
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        type query string false "Transaction type"
// @Param        start_date query string false "Start date"
// @Param        end_date query string false "End date"
// @Param        account_id query uint64 false "Account ID"
// @Param        category_id query uint64 false "Category ID"
// @Param        tag_id query uint64 false "Tag ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions [get]
// @Security     BearerAuth
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

// @Summary      Update transaction
// @Description  Update transaction by ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Param        body body request.UpdateTransactionReq true "update transaction request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id} [put]
// @Security     BearerAuth
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

// @Summary      Delete transaction
// @Description  Delete transaction by ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id} [delete]
// @Security     BearerAuth
func (ctrl *TransactionController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.txnService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

// @Summary      Search transactions
// @Description  Search transactions by keyword with pagination
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        keyword query string false "Search keyword"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/search [get]
// @Security     BearerAuth
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
// @Summary      Split transaction
// @Description  Split a transaction into multiple sub-transactions
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Param        body body request.SplitTransactionReq true "split transaction request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id}/split [post]
// @Security     BearerAuth
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
// @Summary      Get transaction splits
// @Description  Get all sub-transactions of a split transaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id}/splits [get]
// @Security     BearerAuth
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
// @Summary      Merge transaction splits
// @Description  Merge all sub-transactions back into the parent transaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/{id}/merge [post]
// @Security     BearerAuth
func (ctrl *TransactionController) MergeSplits(c *gin.Context) {
	userID := c.GetUint64("user_id")
	parentID := parseIDParam(c, "id")

	if err := ctrl.txnService.MergeSplits(userID, parentID); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
