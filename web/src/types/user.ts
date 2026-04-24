// 用户管理相关类型定义 - 用户接口、角色和列表类型
export interface User {
  id: string
  email: string
  nickname: string
  role: string
  language: string
  timezone: string
  locked: boolean
  created_at: string
  updated_at: string
}

export interface UpdateUserReq {
  nickname?: string
  role?: string
  language?: string
  timezone?: string
}

export interface UserListReq {
  search?: string
  page?: number
  page_size?: number
}
