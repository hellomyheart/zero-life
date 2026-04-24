// 分类相关类型定义 - 分类接口、增删改查
export interface Category {
  id: string
  name: string
  parent_id: string | null
  children: Category[]
  created_at: string
  updated_at: string
}

export interface CreateCategoryReq {
  name: string
  parent_id?: string | null
}

export interface UpdateCategoryReq {
  name?: string
  parent_id?: string | null
}
