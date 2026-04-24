export enum RecurrenceType {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Yearly = 'yearly',
}

export interface RecurringTransaction {
  id: string
  title: string
  type: string
  amount: string
  source_account_id: string
  destination_account_id: string
  category_id: string
  recurrence_type: RecurrenceType
  repeat_interval: number
  start_date: string
  next_date: string
  end_date: string
  active: boolean
  max_repetitions: number
  description: string
  created_at: string
  updated_at: string
}

export interface CreateRecurringTransactionReq {
  title: string
  type: string
  amount: string
  source_account_id: string
  destination_account_id?: string
  category_id?: string
  recurrence_type: RecurrenceType
  repeat_interval: number
  start_date: string
  end_date?: string
  max_repetitions?: number
  description?: string
}

export interface UpdateRecurringTransactionReq {
  title?: string
  amount?: string
  recurrence_type?: RecurrenceType
  repeat_interval?: number
  start_date?: string
  end_date?: string
  active?: boolean
  max_repetitions?: number
  description?: string
}

export interface RecurringTransactionListReq {
  active?: boolean
  page?: number
  page_size?: number
}
