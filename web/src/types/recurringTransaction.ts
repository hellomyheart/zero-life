/**
 * 循环交易相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 重复类型枚举（对应后端 oneof 验证）
 */
export enum RecurrenceType {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Yearly = 'yearly',
}

/**
 * 循环交易信息接口（对应后端 RecurringTransactionResp）
 */
export interface RecurringTransaction {
  id: number
  description: string
  amount: string
  source_id: number
  destination_id: number | null
  category_id: number | null
  notes: string
  recurrence_type: RecurrenceType
  repeat_every: number
  start_date: string
  end_date: string | null
  next_occurrence: string
  is_active: boolean
  created_at: string
  updated_at: string
}

/**
 * 创建循环交易请求接口（对应后端 CreateRecurringTransactionReq）
 */
export interface CreateRecurringTransactionReq {
  description: string
  amount: string
  source_id: number
  destination_id?: number | null
  category_id?: number | null
  notes?: string
  recurrence_type: RecurrenceType
  repeat_every?: number
  start_date: string
  end_date?: string | null
}

/**
 * 更新循环交易请求接口（对应后端 UpdateRecurringTransactionReq）
 */
export interface UpdateRecurringTransactionReq {
  description?: string
  amount?: string
  source_id?: number | null
  destination_id?: number | null
  category_id?: number | null
  notes?: string
  recurrence_type?: RecurrenceType
  repeat_every?: number | null
  start_date?: string
  end_date?: string | null
  is_active?: boolean | null
}

/**
 * 循环交易列表查询请求接口
 */
export interface RecurringTransactionListReq {
  active?: boolean | null
  page?: number
  page_size?: number
}
