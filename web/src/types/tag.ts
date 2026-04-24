// 标签相关类型定义 - 标签接口、增删改查
export interface Tag {
  id: string
  name: string
  color: string
  transaction_count: number
  created_at: string
  updated_at: string
}

export interface CreateTagReq {
  name: string
  color: string
}

export interface UpdateTagReq {
  name?: string
  color?: string
}
