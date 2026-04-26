/**
 * 报表相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 报表查询请求接口
 */
export interface ReportReq {
  start_date: string
  end_date: string
  currency?: string
  account_ids?: number[]
  category_ids?: number[]
}

/**
 * 时段明细响应接口（对应后端 PeriodDetailResp）
 */
export interface PeriodDetailResp {
  period: string
  income: string
  expense: string
  net: string
}

/**
 * 收支报表响应接口（对应后端 IncomeExpenseResp）
 * 后端返回 details（非 income_by_account/expense_by_account）
 */
export interface IncomeExpenseResp {
  total_income: string
  total_expense: string
  net_income: string
  prev_income: string
  prev_expense: string
  prev_net: string
  details: PeriodDetailResp[]
}

/**
 * 分类项响应接口（对应后端 CategoryItemResp）
 */
export interface CategoryItemResp {
  category_id: number
  category_name: string
  amount: string
  percentage: number
}

/**
 * 分类报表响应接口（对应后端 CategoryReportResp）
 * 后端字段名为 income_by_category/expense_by_category
 */
export interface CategoryReportResp {
  income_by_category: CategoryItemResp[]
  expense_by_category: CategoryItemResp[]
}

/**
 * 预算报表项响应接口（对应后端 BudgetReportItemResp）
 */
export interface BudgetReportItemResp {
  budget_id: number
  budget_name: string
  amount: string
  spent: string
  remaining: string
  usage_rate: number
}

/**
 * 预算报表响应接口（对应后端 BudgetReportResp）
 * 后端字段名为 items（非 budgets）
 */
export interface BudgetReportResp {
  items: BudgetReportItemResp[]
}

/**
 * 净资产趋势点响应接口（对应后端 NetWorthPointResp）
 */
export interface NetWorthPointResp {
  date: string
  net_worth: string
}

/**
 * 净值报表响应接口（对应后端 NetWorthResp）
 * 后端仅有 total_net_worth 和 trend，无 total_assets/total_liabilities
 */
export interface NetWorthResp {
  total_net_worth: string
  trend: NetWorthPointResp[]
}

/**
 * 趋势报表项响应接口（对应后端 TrendItemResp）
 * 后端字段名为 period（非 date），无 net 字段
 */
export interface TrendItemResp {
  period: string
  income: string
  expense: string
}

/**
 * 趋势报表响应接口（对应后端 TrendResp）
 */
export interface TrendResp {
  items: TrendItemResp[]
}

/**
 * 标签报表项响应接口（对应后端 TagReportItemResp）
 */
export interface TagItemResp {
  tag_id: number
  tag_name: string
  income: string
  expense: string
}

/**
 * 标签报表响应接口（对应后端 TagReportResp）
 */
export interface TagReportResp {
  items: TagItemResp[]
}
