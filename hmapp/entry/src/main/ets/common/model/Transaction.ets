export class Transaction {
  id: string
  user_id: string
  category_id: string
  amount: number
  transaction_type: string
  account: string
  date: string
  created_at?: number
  updated_at?: number
  deleted_at?: number

  constructor(id: string, user_id: string, category_id: string, amount: number, transaction_type: string, account: string, date: string, created_at?: number, updated_at?: number, deleted_at?: number) {
    this.id = id
    this.user_id = user_id
    this.category_id = category_id
    this.amount = amount
    this.transaction_type = transaction_type
    this.account = account
    this.date = date
    this.created_at = created_at
    this.updated_at = updated_at
    this.deleted_at = deleted_at
  }
}