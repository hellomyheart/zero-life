import type { Account } from '@/types/account'

export interface PiggyBank {
  id: number
  name: string
  account_id: number
  account?: Account
  target_amount: string
  current_amount: string
  target_date: string | null
  notes: string
  percentage: number
  available_deposit: string
  created_at: string
  updated_at: string
}

export interface PiggyEvent {
  id: number
  piggy_bank_id: number
  amount: string
  transaction_id: number | null
  note: string
  created_at: string
}

export interface CreatePiggyBankReq {
  name: string
  account_id: number
  target_amount: string
  target_date?: string | null
  notes?: string
}

export interface UpdatePiggyBankReq {
  name?: string
  target_amount?: string
  target_date?: string | null
  notes?: string
  clear_notes?: boolean
}

export interface AddAmountReq {
  amount: string
  note?: string
}

export interface RemoveAmountReq {
  amount: string
  note?: string
}
