export enum TransactionType {
  Deposit = 'deposit',
  Withdrawal = 'withdrawal',
  Transfer = 'transfer',
}

export interface Split {
  id: string
  amount: string
  category_id: string
  category_name: string
  tag_ids: string[]
  tag_names: string[]
  description: string
}

export interface Transaction {
  id: string
  type: TransactionType
  date: string
  description: string
  amount: string
  source_account_id: string
  source_account_name: string
  destination_account_id: string
  destination_account_name: string
  category_id: string
  category_name: string
  tag_ids: string[]
  tag_names: string[]
  splits: Split[]
  created_at: string
  updated_at: string
}

export interface CreateSplitReq {
  amount: string
  category_id: string
  tag_ids: string[]
  description: string
}

export interface CreateTransactionReq {
  type: TransactionType
  date: string
  description: string
  amount: string
  source_account_id: string
  destination_account_id?: string
  category_id?: string
  tag_ids?: string[]
  splits?: CreateSplitReq[]
}

export interface UpdateTransactionReq {
  type?: TransactionType
  date?: string
  description?: string
  amount?: string
  source_account_id?: string
  destination_account_id?: string
  category_id?: string
  tag_ids?: string[]
  splits?: CreateSplitReq[]
}

export interface TransactionListReq {
  type?: TransactionType
  start_date?: string
  end_date?: string
  source_account_id?: string
  destination_account_id?: string
  category_id?: string
  tag_ids?: string[]
  keyword?: string
  page?: number
  page_size?: number
}
