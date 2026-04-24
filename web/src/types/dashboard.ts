// 仪表盘相关类型定义 - 首页汇总数据接口
export interface BudgetAlert {
  budget_id: string
  budget_name: string
  amount: string
  spent: string
  usage_rate: number
  status: string
}

export interface BillReminder {
  bill_id: string
  bill_name: string
  amount: string
  next_due_date: string
  is_overdue: boolean
}

export interface DashboardResp {
  total_income: string
  total_expense: string
  net_income: string
  total_balance: string
  budget_alerts: BudgetAlert[]
  bill_reminders: BillReminder[]
  recent_transactions: {
    id: string
    date: string
    description: string
    amount: string
    type: string
  }[]
}
