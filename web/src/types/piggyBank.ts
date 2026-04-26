/**
 * 存钱罐相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 存钱罐信息接口（对应后端 PiggyBankResp）
 * 后端不返回 start_date 和 order 字段
 */
export interface PiggyBank {
  id: number
  name: string
  account_id: number
  target_amount: string
  current_amount: string
  target_date: string | null
  notes: string
  percentage: number
  created_at: string
  updated_at: string
}

/**
 * 存钱罐事件接口（对应后端 PiggyEventResp）
 * 后端不返回 type 字段，有 transaction_id 字段
 */
export interface PiggyEvent {
  id: number
  piggy_bank_id: number
  amount: string
  transaction_id: number | null
  note: string
  created_at: string
}

/**
 * 创建存钱罐请求接口（对应后端 CreatePiggyBankReq）
 */
export interface CreatePiggyBankReq {
  name: string
  account_id: number
  target_amount: string
  target_date?: string | null
  notes?: string
}

/**
 * 更新存钱罐请求接口（对应后端 UpdatePiggyBankReq）
 */
export interface UpdatePiggyBankReq {
  name?: string
  target_amount?: string
  target_date?: string | null
  notes?: string
}

/**
 * 存取款请求接口（对应后端 AddAmountReq）
 */
export interface AddAmountReq {
  amount: string
  note?: string
}

/**
 * 取款请求接口（对应后端 RemoveAmountReq）
 */
export interface RemoveAmountReq {
  amount: string
  note?: string
}

