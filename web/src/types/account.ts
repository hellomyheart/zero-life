export enum AccountType {
  Asset = 'asset',
  Expense = 'expense',
  Income = 'income',
  Liability = 'liability',
}

export interface Account {
  id: string
  name: string
  type: AccountType
  currency: string
  initial_balance: string
  balance: string
  is_virtual: boolean
  created_at: string
  updated_at: string
}

export interface CreateAccountReq {
  name: string
  type: AccountType
  currency: string
  initial_balance: string
  is_virtual?: boolean
}

export interface UpdateAccountReq {
  name?: string
  type?: AccountType
  currency?: string
  initial_balance?: string
  is_virtual?: boolean
}

export interface AccountListReq {
  type?: AccountType
  currency?: string
  page?: number
  page_size?: number
}
