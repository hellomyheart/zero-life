export class User {
  id: string
  nickname?: string
  email?: string
  phone?: string
  password?: string
  avatar?: string
  created_at?: number
  updated_at?: number
  deleted_at?: number

  constructor(id: string, nickname?: string, email?: string, phone?: string, password?: string, avatar?: string, created_at?: number, updated_at?: number, deleted_at?: number) {
    this.id = id
    this.nickname = nickname
    this.email = email
    this.phone = phone
    this.password = password
    this.avatar = avatar
    this.created_at = created_at
    this.updated_at = updated_at
    this.deleted_at = deleted_at
  }
}