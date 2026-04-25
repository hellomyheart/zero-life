// Package controller 提供HTTP请求处理控制器
// BillController 账单控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// BillController 账单管理控制器
// 处理账单的创建、查询、更新和删除等HTTP请求
// 账单用于记录周期性支出（如房租、水电费等），支持日/周/月/年重复规则
type BillController struct {
	billService *service.BillService // 账单业务服务
}

// NewBillController 创建账单控制器实例
// 参数：
//   - billService: 账单业务服务实例
// 返回：
//   - *BillController: 账单控制器实例
func NewBillController(billService *service.BillService) *BillController {
	return &BillController{billService: billService}
}

// Create 创建账单
// 接收创建账单请求，设置账单名称、金额、重复规则和下次到期日
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateBillReq（包含名称、金额、重复规则、下次到期日等）
// 响应：创建成功的账单信息
// 业务规则：账单金额必须大于0，重复规则仅支持daily/weekly/monthly/yearly
// @Summary      Create bill
// @Description  Create a new bill for the current user
// @Tags         bills
// @Accept       json
// @Produce      json
// @Param        body body request.CreateBillReq true "create bill request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/bills [post]
// @Security     BearerAuth
func (ctrl *BillController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.billService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取账单详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（账单ID）
// 响应：账单详细信息
// @Summary      Get bill
// @Description  Get bill details by ID
// @Tags         bills
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Bill ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/bills/{id} [get]
// @Security     BearerAuth
func (ctrl *BillController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.billService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取账单列表
// 返回当前用户的所有账单
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：账单列表
// @Summary      List bills
// @Description  Get all bills for the current user
// @Tags         bills
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/bills [get]
// @Security     BearerAuth
func (ctrl *BillController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.billService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新账单
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（账单ID）
// 请求体：UpdateBillReq（包含需要更新的字段）
// 响应：更新后的账单信息
// @Summary      Update bill
// @Description  Update bill by ID
// @Tags         bills
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Bill ID"
// @Param        body body request.UpdateBillReq true "update bill request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/bills/{id} [put]
// @Security     BearerAuth
func (ctrl *BillController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.billService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除账单
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（账单ID）
// 响应：删除成功返回nil
// @Summary      Delete bill
// @Description  Delete bill by ID
// @Tags         bills
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Bill ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/bills/{id} [delete]
// @Security     BearerAuth
func (ctrl *BillController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.billService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
