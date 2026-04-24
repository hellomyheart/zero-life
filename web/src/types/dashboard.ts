/**
 * 仪表盘相关类型定义
 * 定义首页仪表盘展示所需的汇总数据结构
 */

/**
 * 预算预警接口
 * 当预算使用率接近或超过上限时产生的预警信息
 */
export interface BudgetAlert {
  budget_id: string            // 预算ID
  budget_name: string          // 预算名称
  amount: string               // 预算金额上限
  spent: string                // 已花费金额
  usage_rate: number           // 使用率（0-1之间，超过1表示超支）
  status: string               // 预警状态，如"warning"、"exceeded"等
}

/**
 * 账单提醒接口
 * 即将到期或已逾期的账单提醒信息
 */
export interface BillReminder {
  bill_id: string              // 账单ID
  bill_name: string            // 账单名称
  amount: string               // 账单金额
  next_due_date: string        // 下次到期日期
  is_overdue: boolean          // 是否已逾期
}

/**
 * 仪表盘响应接口
 * 首页仪表盘展示的所有汇总数据
 */
export interface DashboardResp {
  total_income: string         // 总收入
  total_expense: string        // 总支出
  net_income: string           // 净收入（收入 - 支出）
  total_balance: string        // 总余额（所有资产账户余额之和）
  budget_alerts: BudgetAlert[] // 预算预警列表
  bill_reminders: BillReminder[] // 账单提醒列表
  recent_transactions: {       // 最近交易列表
    id: string                 // 交易ID
    date: string               // 交易日期
    description: string        // 交易描述
    amount: string             // 交易金额
    type: string               // 交易类型（deposit/withdrawal/transfer）
  }[]
}
