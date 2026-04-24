// 报表相关类型定义 - 各类报表数据接口
export interface ReportReq {
  start_date: string
  end_date: string
  currency?: string
  account_ids?: string[]
  category_ids?: string[]
}

export interface IncomeExpenseResp {
  total_income: string
  total_expense: string
  net_income: string
  income_by_account: { account_id: string; account_name: string; amount: string }[]
  expense_by_account: { account_id: string; account_name: string; amount: string }[]
}

export interface CategoryItemResp {
  category_id: string
  category_name: string
  amount: string
  percentage: number
}

export interface CategoryReportResp {
  income_categories: CategoryItemResp[]
  expense_categories: CategoryItemResp[]
}

export interface BudgetReportResp {
  budgets: {
    budget_id: string
    budget_name: string
    amount: string
    spent: string
    usage_rate: number
    status: string
  }[]
}

export interface NetWorthResp {
  items: { date: string; net_worth: string }[]
  total_assets: string
  total_liabilities: string
  net_worth: string
}

export interface TrendResp {
  items: { date: string; income: string; expense: string; net: string }[]
}
