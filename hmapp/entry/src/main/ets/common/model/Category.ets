export class Category {
  id: string
  name: string
  transaction_type: string
  created_at?: number
  updated_at?: number
  deleted_at?: number

  constructor(id: string, name: string, transaction_type: string, created_at?: number, updated_at?: number, deleted_at?: number) {
    this.id = id
    this.name = name
    this.transaction_type = transaction_type
    this.created_at = created_at
    this.updated_at = updated_at
    this.deleted_at = deleted_at
  }
}