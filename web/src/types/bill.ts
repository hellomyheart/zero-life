export enum RepeatRule {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Quarterly = 'quarterly',
  Yearly = 'yearly',
}

export interface Bill {
  id: string
  name: string
  amount: string
  account_id: string
  account_name: string
  category_id: string
  category_name: string
  repeat_rule: RepeatRule
  next_due_date: string
  is_overdue: boolean
  description: string
  created_at: string
  updated_at: string
}

export interface CreateBillReq {
  name: string
  amount: string
  account_id: string
  category_id: string
  repeat_rule: RepeatRule
  next_due_date: string
  description?: string
}

export interface UpdateBillReq {
  name?: string
  amount?: string
  account_id?: string
  category_id?: string
  repeat_rule?: RepeatRule
  next_due_date?: string
  description?: string
}
