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
