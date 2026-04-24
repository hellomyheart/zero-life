import { get } from '@/utils/request'
import type { IncomeExpenseResp, CategoryReportResp, BudgetReportResp, NetWorthResp, TrendResp, ReportReq } from '@/types/report'

export function incomeExpense(params: ReportReq) {
  return get<IncomeExpenseResp>('/reports/income-expense', params as unknown as Record<string, unknown>)
}

export function category(params: ReportReq) {
  return get<CategoryReportResp>('/reports/category', params as unknown as Record<string, unknown>)
}

export function budget(params: ReportReq) {
  return get<BudgetReportResp>('/reports/budget', params as unknown as Record<string, unknown>)
}

export function netWorth(params: ReportReq) {
  return get<NetWorthResp>('/reports/net-worth', params as unknown as Record<string, unknown>)
}

export function trend(params: ReportReq) {
  return get<TrendResp>('/reports/trend', params as unknown as Record<string, unknown>)
}
