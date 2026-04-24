// 通用类型定义 - API响应、分页、筛选分页
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface PageParams {
  page?: number
  page_size?: number
}
