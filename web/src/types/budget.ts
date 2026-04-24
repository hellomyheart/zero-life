// 预算相关类型定义 - 预算接口、增删改查和历史记录
export enum BudgetPeriod {
  Monthly = 'monthly',
  Quarterly = 'quarterly',
  Yearly = 'yearly',
}

export interface Budget {
  id: string
  name: string
  amount: string
  spent: string
  period: BudgetPeriod
  category_ids: string[]
  category_names: string[]
  start_date: string
  status: string
  created_at: string
  updated_at: string
}

export interface CreateBudgetReq {
  name: string
  amount: string
  period: BudgetPeriod
  category_ids: string[]
  start_date: string
}

export interface UpdateBudgetReq {
  name?: string
  amount?: string
  period?: BudgetPeriod
  category_ids?: string[]
  start_date?: string
}

export interface BudgetHistory {
  period: string
  amount: string
  spent: string
  usage_rate: number
}
