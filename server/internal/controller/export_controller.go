// Package controller 提供HTTP请求处理控制器
// ExportController 数据导出控制器
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ExportController 数据导出控制器
// 处理交易、账户、预算、分类、标签等数据的导出请求
// 支持CSV和JSON两种导出格式，以文件下载方式返回
type ExportController struct {
	service *service.ExportService // 导出业务服务
}

// NewExportController 创建导出控制器实例
// 参数：
//   - service: 导出业务服务实例
// 返回：
//   - *ExportController: 导出控制器实例
func NewExportController(service *service.ExportService) *ExportController {
	return &ExportController{service: service}
}

// ExportTransactions 导出交易数据
// 根据日期范围和格式导出交易记录
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：ExportReq（包含start_date、end_date、format）
// 响应：文件流下载（CSV或JSON格式）
func (c *ExportController) ExportTransactions(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportTransactions(userID, req.StartDate, req.EndDate, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

// ExportAccounts 导出账户数据
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：ExportReq（包含format）
// 响应：文件流下载（CSV或JSON格式）
func (c *ExportController) ExportAccounts(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportAccounts(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

// ExportBudgets 导出预算数据
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：ExportReq（包含format）
// 响应：文件流下载（CSV或JSON格式）
func (c *ExportController) ExportBudgets(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportBudgets(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

// ExportCategories 导出分类数据
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：ExportReq（包含format）
// 响应：文件流下载（CSV或JSON格式）
func (c *ExportController) ExportCategories(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportCategories(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

// ExportTags 导出标签数据
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：ExportReq（包含format）
// 响应：文件流下载（CSV或JSON格式）
func (c *ExportController) ExportTags(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportTags(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}
