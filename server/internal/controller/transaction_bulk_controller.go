// Package controller 提供HTTP请求处理控制器
// TransactionBulkController 交易批量操作控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// TransactionBulkController 交易批量操作控制器
// 处理交易的批量编辑、批量删除、类型转换和克隆等HTTP请求
type TransactionBulkController struct {
	bulkService *service.TransactionBulkService // 批量操作业务服务
}

// NewTransactionBulkController 创建批量操作控制器实例
// 参数：
//   - bulkService: 批量操作业务服务实例
// 返回：
//   - *TransactionBulkController: 批量操作控制器实例
func NewTransactionBulkController(bulkService *service.TransactionBulkService) *TransactionBulkController {
	return &TransactionBulkController{bulkService: bulkService}
}

// BulkEdit 批量编辑交易
// 批量修改指定交易的分类、备注和标签
// 参数：
//   - c: Gin上下文，包含用户身份
// 请求体：BulkEditReq（包含交易ID列表、分类ID、备注、标签ID列表）
// 响应：编辑成功返回nil
// @Summary      Bulk edit transactions
// @Description  Batch modify category, notes, and tags for selected transactions
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body body request.BulkEditReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/bulk/edit [post]
// @Security     BearerAuth
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

// BulkDelete 批量删除交易
// 批量删除指定ID列表中的交易
// 参数：
//   - c: Gin上下文，包含用户身份
// 请求体：BulkDeleteReq（包含交易ID列表）
// 响应：删除成功返回nil
// @Summary      Bulk delete transactions
// @Description  Batch delete transactions by ID list
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body body request.BulkDeleteReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/bulk/delete [post]
// @Security     BearerAuth
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

// ConvertType 转换交易类型
// 将指定交易从一种类型转换为另一种类型（如支出转收入）
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（交易ID）
// 请求体：ConvertReq（包含目标类型、源账户ID、目标账户ID）
// 响应：转换后的交易信息
// @Summary      Convert transaction type
// @Description  Convert a transaction from one type to another
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Param        body body request.ConvertReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/bulk/{id}/convert [post]
// @Security     BearerAuth
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

// Clone 克隆交易
// 复制指定交易的所有信息创建一笔新交易
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（原交易ID）
// 响应：克隆后的新交易信息
// @Summary      Clone transaction
// @Description  Clone a transaction creating a new copy with all the same info
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transactions/bulk/{id}/clone [post]
// @Security     BearerAuth
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
