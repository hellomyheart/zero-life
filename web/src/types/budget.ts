/**
 * 预算相关类型定义
 * 预算用于设定某个时间段内特定分类的支出上限，帮助控制开支
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 预算周期枚举
 * - Monthly: 月度预算，每月重置
 * - Yearly: 年度预算，每年重置
 */
export enum BudgetPeriod {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Quarterly = 'quarterly',
  Yearly = 'yearly',
}

/**
 * 预算中的分类信息（后端 CategoryResp 嵌套）
 */
export interface BudgetCategoryResp {
  id: number                    // 分类ID（后端 uint64）
  name: string                  // 分类名称
  parent_id: number | null      // 父分类ID
  icon: string                  // 图标
  notes: string                 // 备注
  sort_order: number            // 排序
  created_at: string            // 创建时间
  updated_at: string            // 更新时间
}

/**
 * 预算信息接口（对应后端 BudgetResp）
 * 包含预算基本信息和执行情况（已支出、剩余、使用率、状态）
 */
export interface Budget {
  id: number                    // 预算唯一标识（后端 uint64）
  name: string                  // 预算名称，如"餐饮月度预算"
  amount: string                // 预算金额上限
  period: BudgetPeriod          // 预算周期：monthly/yearly
  is_enabled: boolean           // 是否启用
  categories: BudgetCategoryResp[] // 关联分类列表（对象数组）
  spent: string                 // 当前周期已花费金额
  remaining: string             // 当前周期剩余金额
  usage_rate: number            // 使用率（0~1+，1=100%，>1=超支）
  status: string                // 状态：normal/warning/overspent
  created_at: string            // 创建时间
  updated_at: string            // 最后更新时间
}

/**
 * 创建预算请求接口（对应后端 CreateBudgetReq）
 */
export interface CreateBudgetReq {
  name: string                  // 预算名称（必填）
  amount: string                // 预算金额上限（必填，必须大于0）
  period: BudgetPeriod          // 预算周期（必填，monthly/yearly）
  category_ids: number[]        // 关联分类ID列表（必填，至少1个）
}

/**
 * 更新预算请求接口（对应后端 UpdateBudgetReq）
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateBudgetReq {
  name?: string                 // 预算名称
  amount?: string               // 预算金额上限
  period?: BudgetPeriod         // 预算周期
  category_ids?: number[]       // 关联分类ID列表
  is_enabled?: boolean          // 是否启用
}

/**
 * 预算历史记录接口（对应后端 BudgetHistoryResp）
 * 记录预算在每个周期内的使用情况
 */
export interface BudgetHistory {
  id: number                    // 历史记录ID（后端 uint64）
  period_start: string          // 周期开始时间
  period_end: string            // 周期结束时间
  amount: string                // 该周期的预算金额
  spent: string                 // 该周期的实际支出
  usage_rate: number            // 使用率（0~1+）
  created_at: string            // 记录创建时间
}
