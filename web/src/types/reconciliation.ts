/**
 * 对账相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 对账记录接口（对应后端 ReconciliationResp）
 */
export interface Reconciliation {
  id: number
  account_id: number
  start_date: string
  end_date: string
  starting_balance: string
  ending_balance: string
  book_balance: string
  difference: string
  status: string
  created_at: string
  updated_at: string
}

/**
 * 创建对账记录请求接口（对应后端 CreateReconciliationReq）
 * 后端字段名为 starting_balance/ending_balance（非 start_balance/end_balance）
 * 后端不接受 notes 字段
 */
export interface CreateReconciliationReq {
  account_id: number
  start_date: string
  end_date: string
  starting_balance: string
  ending_balance: string
}

/**
 * 更新对账记录请求接口（对应后端 UpdateReconciliationReq）
 */
export interface UpdateReconciliationReq {
  ending_balance?: string
  status?: string
}

/**
 * 对账记录列表查询请求接口
 */
export interface ReconciliationListReq {
  account_id?: number | null
  page?: number
  page_size?: number
}
