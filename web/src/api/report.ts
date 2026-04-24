/**
 * 报表相关API接口
 * 提供各类财务报表数据的查询功能，包括收支报表、分类报表、标签报表、
 * 预算报表、净值报表和趋势报表
 */
import { get } from '@/utils/request'
import type { IncomeExpenseResp, CategoryReportResp, BudgetReportResp, NetWorthResp, TrendResp, ReportReq } from '@/types/report'

/**
 * 获取收支报表
 * 返回指定时间段内的总收入、总支出、净收入，以及按账户分组的收支明细
 * @param params - 报表查询参数（开始日期、结束日期、可选的货币和账户/分类筛选）
 * @returns 收支报表数据
 * @endpoint GET /reports/income-expense
 */
export function incomeExpense(params: ReportReq) {
  return get<IncomeExpenseResp>('/reports/income-expense', params as unknown as Record<string, unknown>)
}

/**
 * 获取分类报表
 * 返回按分类汇总的收入和支出数据，包含每个分类的金额和占比
 * @param params - 报表查询参数
 * @returns 分类报表数据
 * @endpoint GET /reports/category
 */
export function category(params: ReportReq) {
  return get<CategoryReportResp>('/reports/category', params as unknown as Record<string, unknown>)
}

/**
 * 获取标签报表
 * 返回按标签汇总的收入和支出数据，包含每个标签的金额和占比
 * @param params - 报表查询参数
 * @returns 标签报表数据
 * @endpoint GET /reports/tag
 */
export function tag(params: ReportReq) {
  return get<CategoryReportResp>('/reports/tag', params as unknown as Record<string, unknown>)
}

/**
 * 获取预算报表
 * 返回各预算的使用情况，包含预算金额、已花费金额、使用率和状态
 * @param params - 报表查询参数
 * @returns 预算报表数据
 * @endpoint GET /reports/budget
 */
export function budget(params: ReportReq) {
  return get<BudgetReportResp>('/reports/budget', params as unknown as Record<string, unknown>)
}

/**
 * 获取净值报表
 * 返回总资产、总负债和净值，以及按日期的净值变化趋势
 * @param params - 报表查询参数
 * @returns 净值报表数据
 * @endpoint GET /reports/net-worth
 */
export function netWorth(params: ReportReq) {
  return get<NetWorthResp>('/reports/net-worth', params as unknown as Record<string, unknown>)
}

/**
 * 获取趋势报表
 * 返回按日期汇总的收入、支出和净额变化趋势，用于绘制趋势图
 * @param params - 报表查询参数
 * @returns 趋势报表数据
 * @endpoint GET /reports/trend
 */
export function trend(params: ReportReq) {
  return get<TrendResp>('/reports/trend', params as unknown as Record<string, unknown>)
}
