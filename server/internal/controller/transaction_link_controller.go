// Package controller 提供HTTP请求处理控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// TransactionLinkController 交易关联控制器
// 负责处理交易关联相关的HTTP请求，包括交易关联的创建、查询和删除
// 所有操作均需传入用户ID，确保用户只能访问自己交易的关联数据
type TransactionLinkController struct {
	service *service.TransactionLinkService
}

// NewTransactionLinkController 创建交易关联控制器实例
// 参数：
//   - service: 交易关联业务服务实例
// 返回：
//   - *TransactionLinkController: 交易关联控制器实例
func NewTransactionLinkController(service *service.TransactionLinkService) *TransactionLinkController {
	return &TransactionLinkController{service: service}
}

// Create 创建交易关联
// 从上下文中提取用户ID，绑定请求参数，调用服务层创建交易关联
// 请求方法：POST
// 请求路径：/transaction-links
// 请求体：CreateTransactionLinkReq（JSON格式）
// @Summary      Create transaction link
// @Description  Create a link between two transactions
// @Tags         transaction-links
// @Accept       json
// @Produce      json
// @Param        body body request.CreateTransactionLinkReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transaction-links [post]
// @Security     BearerAuth
func (ctrl *TransactionLinkController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateTransactionLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取交易关联列表（分页）
// 从上下文中提取用户ID，确保只能查询属于当前用户交易的关联记录
// 请求方法：GET
// 请求路径：/transaction-links
// 查询参数：TransactionLinkListReq
// @Summary      List transaction links
// @Description  Get paginated list of transaction links
// @Tags         transaction-links
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transaction-links [get]
// @Security     BearerAuth
func (ctrl *TransactionLinkController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.TransactionLinkListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Delete 删除交易关联
// 从上下文中提取用户ID，确保只能删除属于当前用户交易的关联记录
// 请求方法：DELETE
// 请求路径：/transaction-links/:id
// @Summary      Delete transaction link
// @Description  Delete a transaction link by ID
// @Tags         transaction-links
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Transaction link ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/transaction-links/{id} [delete]
// @Security     BearerAuth
func (ctrl *TransactionLinkController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
