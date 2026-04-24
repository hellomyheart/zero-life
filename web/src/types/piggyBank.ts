// 存钱罐相关类型定义 - 存钱罐接口、存取款和事件记录
export interface PiggyBank {
  id: string
  name: string
  account_id: string
  target_amount: string
  current_amount: string
  start_date: string
  target_date: string
  order: number
  notes: string
  percentage: number
  created_at: string
  updated_at: string
}

export interface PiggyEvent {
  id: string
  piggy_bank_id: string
  amount: string
  type: string
  note: string
  created_at: string
}

export interface CreatePiggyBankReq {
  name: string
  account_id: string
  target_amount: string
  target_date?: string
  notes?: string
}

export interface UpdatePiggyBankReq {
  name?: string
  target_amount?: string
  target_date?: string
  notes?: string
}

export interface AddAmountReq {
  amount: string
  note?: string
}

export interface PiggyBankListReq {
  page?: number
  page_size?: number
}
