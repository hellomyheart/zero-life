export interface Reconciliation {
  id: string
  account_id: string
  start_date: string
  end_date: string
  start_balance: string
  end_balance: string
  book_balance: string
  difference: string
  status: string
  notes: string
  created_at: string
  updated_at: string
}

export interface CreateReconciliationReq {
  account_id: string
  start_date: string
  end_date: string
  start_balance: string
  end_balance: string
  notes?: string
}

export interface UpdateReconciliationReq {
  end_balance?: string
  status?: string
  notes?: string
}

export interface ReconciliationListReq {
  account_id?: string
  page?: number
  page_size?: number
}
