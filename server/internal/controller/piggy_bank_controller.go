// Package controller 提供HTTP请求处理控制器
// PiggyBankController 存钱罐控制器
// 负责处理存钱罐相关的所有HTTP请求，包括创建、查询、更新、删除、
// 存取金额、事件查询、排序和重置历史等操作。
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// PiggyBankController 存钱罐控制器
// 封装了存钱罐相关的业务逻辑调用，将HTTP请求转发给service层处理。
// 每个方法对应一个API端点，负责参数解析、请求校验和响应返回。
type PiggyBankController struct {
	service *service.PiggyBankService
}

// NewPiggyBankController 创建存钱罐控制器实例
// 参数：
//   - service: 存钱罐服务实例，负责具体的业务逻辑处理
// 返回值：
//   - *PiggyBankController: 初始化后的控制器实例
func NewPiggyBankController(service *service.PiggyBankService) *PiggyBankController {
	return &PiggyBankController{service: service}
}

// Create 创建存钱罐
// 从请求体中解析存钱罐名称、目标金额等信息，调用service层创建新的存钱罐。
// 参数：
//   - ctx: Gin上下文，包含请求信息和用户身份（通过中间件注入的user_id）
// 请求体：CreatePiggyBankReq（包含名称、目标金额等字段）
// 响应：创建成功的存钱罐信息
func (c *PiggyBankController) Create(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.CreatePiggyBankReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	resp, err := c.service.Create(userID, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// Get 获取单个存钱罐详情
// 根据URL路径中的存钱罐ID，查询该存钱罐的详细信息。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 响应：存钱罐详细信息
func (c *PiggyBankController) Get(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	resp, err := c.service.Get(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// List 获取当前用户的所有存钱罐列表
// 查询当前登录用户名下的所有存钱罐，返回列表信息。
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 响应：存钱罐列表
func (c *PiggyBankController) List(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	resp, err := c.service.List(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// Update 更新存钱罐信息
// 根据URL路径中的存钱罐ID和请求体中的更新字段，修改存钱罐信息。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 请求体：UpdatePiggyBankReq（包含需要更新的字段）
// 响应：更新后的存钱罐信息
func (c *PiggyBankController) Update(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	var req request.UpdatePiggyBankReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	resp, err := c.service.Update(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// Delete 删除存钱罐
// 根据URL路径中的存钱罐ID，删除对应的存钱罐及其关联数据。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 响应：删除成功返回nil
func (c *PiggyBankController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	if err := c.service.Delete(userID, id); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// AddAmount 存入金额
// 向指定存钱罐中增加金额，并记录一条存入事件。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 请求体：AddAmountReq（包含存入金额等字段）
// 响应：更新后的存钱罐信息
func (c *PiggyBankController) AddAmount(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	var req request.AddAmountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	resp, err := c.service.AddAmount(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// RemoveAmount 取出金额
// 从指定存钱罐中减少金额，并记录一条取出事件。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 请求体：RemoveAmountReq（包含取出金额等字段）
// 响应：更新后的存钱罐信息
func (c *PiggyBankController) RemoveAmount(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	var req request.RemoveAmountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	resp, err := c.service.RemoveAmount(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// GetEvents 获取存钱罐的事件记录
// 查询指定存钱罐的所有存取事件历史记录。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 响应：事件记录列表
func (c *PiggyBankController) GetEvents(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	resp, err := c.service.GetEvents(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}

// Reorder 调整存钱罐排序
// 批量修改存钱罐的显示顺序，接收一个map，key为存钱罐ID，value为目标排序位置。
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 请求体：包含orders字段，类型为map[uint64]int，key是存钱罐ID，value是排序序号
// 响应：排序成功返回nil
func (c *PiggyBankController) Reorder(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req struct {
		Orders map[uint64]int `json:"orders" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.service.Reorder(userID, req.Orders); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

// ResetHistory 重置存钱罐历史记录
// 清空指定存钱罐的所有存取事件记录，将当前金额归零。
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（存钱罐ID）
// 响应：重置后的存钱罐信息
func (c *PiggyBankController) ResetHistory(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	resp, err := c.service.ResetHistory(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}
