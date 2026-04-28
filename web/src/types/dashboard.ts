/**
 * 仪表盘相关类型定义
 * 定义首页仪表盘展示所需的汇总数据结构
 * 注意：字段名必须与后端 JSON 响应完全对应
 */

/**
 * 预算预警接口
 * 当预算使用率接近或超过上限时产生的预警信息
 */
export interface BudgetAlert {
  budget_id: number            // 预算ID（后端返回 uint64）
  budget_name: string          // 预算名称
  amount: string               // 预算金额上限
  spent: string                // 已花费金额
  usage_rate: number           // 使用率（0-1之间，超过1表示超支）
  status: string               // 预警状态，如"warning"、"exceeded"等
}

/**
 * 循环交易提醒接口
 * 即将到期的循环交易提醒信息
 */
export interface RecurringReminder {
  recurring_id: number         // 循环交易ID（后端返回 uint64）
  name: string                 // 名称
  amount: string               // 金额
  next_due: string             // 下次到期日期（后端字段名为 next_due）
}

/**
 * 仪表盘响应接口
 * 首页仪表盘展示的所有汇总数据
 * 字段名与后端 DashboardResp JSON tag 完全对应
 */
export interface DashboardResp {
  month_income: string         // 本月收入（后端字段名 month_income）
  month_expense: string        // 本月支出（后端字段名 month_expense）
  net_income: string           // 净收入（收入 - 支出）
  total_balance: string        // 总余额（所有资产账户余额之和）
  budget_alerts: BudgetAlert[] // 预算预警列表
  recurring_reminders: RecurringReminder[] // 循环交易提醒列表
  recent_txns: {               // 最近交易列表（后端字段名 recent_txns）
    id: number                 // 交易ID
    date: string               // 交易日期
    description: string        // 交易描述
    amount: string             // 交易金额
    type: string               // 交易类型（deposit/withdrawal/transfer）
  }[]
}
