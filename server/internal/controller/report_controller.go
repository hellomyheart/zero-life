// Package controller 提供HTTP请求处理控制器
// ReportController 报表控制器
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ReportController 报表管理控制器
// 处理各类报表数据查询请求，包括收支报表、分类报表、预算报表、
// 净资产报表、趋势报表、标签报表和审计报表
type ReportController struct {
	reportService *service.ReportService // 报表业务服务
}

// NewReportController 创建报表控制器实例
// 参数：
//   - reportService: 报表业务服务实例
// 返回：
//   - *ReportController: 报表控制器实例
func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

// IncomeExpense 收支报表
// 查询指定时间段内的收入和支出汇总数据，支持按日/周/月/季/年粒度展示
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date、granularity）
// 响应：IncomeExpenseResp（含总收入、总支出、净收入、分时段明细）
func (ctrl *ReportController) IncomeExpense(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.IncomeExpense(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Category 分类报表
// 查询指定时间段内按分类汇总的支出和收入分布数据
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date）
// 响应：CategoryReportResp（含支出分类分布和收入分类分布）
func (ctrl *ReportController) Category(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Category(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Budget 预算报表
// 查询指定时间段内各预算的执行情况
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date）
// 响应：BudgetReportResp（含各预算的金额、已支出、剩余、使用率）
func (ctrl *ReportController) Budget(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Budget(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// NetWorth 净资产报表
// 查询指定时间段内的净资产变化趋势
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date）
// 响应：NetWorthResp（含总净资产和趋势数据）
func (ctrl *ReportController) NetWorth(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.NetWorth(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Trend 趋势报表
// 查询指定时间段内收支趋势数据
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date、granularity）
// 响应：TrendResp（含各时段的收入和支出数据）
func (ctrl *ReportController) Trend(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Trend(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Tag 标签报表
// 查询指定时间段内按标签汇总的收支数据
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：ReportReq（包含start_date、end_date）
// 响应：TagReportResp（含各标签的收入和支出数据）
func (ctrl *ReportController) Tag(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ReportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reportService.Tag(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Audit 审计报表
// 查询指定账户在指定时间段内的交易审计数据，含每笔交易的余额流水
// 参数：
//   - c: Gin上下文，包含用户身份
// 查询参数：account_id（必填）、start_date、end_date、reconciled（是否已对账过滤）
// 响应：AuditReportResp（含账户信息、期初余额、期末余额、交易明细列表）
func (ctrl *ReportController) Audit(c *gin.Context) {
	userID := c.GetUint64("user_id")

	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var reconciled *bool
	if r := c.Query("reconciled"); r != "" {
		val, err := strconv.ParseBool(r)
		if err == nil {
			reconciled = &val
		}
	}

	result, err := ctrl.reportService.AuditReport(userID, accountID, startDate, endDate, reconciled)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
