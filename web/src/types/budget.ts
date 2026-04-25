/**
 * 预算相关类型定义
 * 预算用于设定某个时间段内特定分类的支出上限，帮助控制开支
 */

/**
 * 预算周期枚举
 * 定义预算的时间范围
 * - Monthly: 月度预算
 * - Quarterly: 季度预算
 * - Yearly: 年度预算
 */
export enum BudgetPeriod {
  Monthly = 'monthly',
  Quarterly = 'quarterly',
  Yearly = 'yearly',
}

/**
 * 预算信息接口
 * 表示一个预算的完整信息，包含已花费金额和使用率
 */
export interface Budget {
  id: number                    // 预算唯一标识
  name: string                  // 预算名称，如"餐饮月度预算"
  amount: string                // 预算金额上限
  spent: string                 // 当前周期已花费金额
  period: BudgetPeriod          // 预算周期
  category_ids: number[]        // 关联的分类ID列表（预算针对这些分类的支出）
  category_names: string[]      // 关联的分类名称列表（冗余字段，方便展示）
  start_date: string           // 预算开始日期
  status: string               // 预算状态，如"active"、"exceeded"等
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建预算请求接口
 */
export interface CreateBudgetReq {
  name: string                 // 预算名称（必填）
  amount: string               // 预算金额上限（必填）
  period: BudgetPeriod         // 预算周期（必填）
  category_ids: number[]        // 关联分类ID列表（必填）
  start_date: string           // 开始日期（必填）
}

/**
 * 更新预算请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateBudgetReq {
  name?: string                // 预算名称
  amount?: string              // 预算金额上限
  period?: BudgetPeriod        // 预算周期
  category_ids?: number[]       // 关联分类ID列表
  start_date?: string          // 开始日期
}

/**
 * 预算历史记录接口
 * 记录预算在每个周期内的使用情况
 */
export interface BudgetHistory {
  period: string               // 周期标识，如"2024-01"
  amount: string               // 该周期的预算金额
  spent: string                // 该周期的实际花费
  usage_rate: number           // 使用率（0-1之间，1表示已用完，超过1表示超支）
}
