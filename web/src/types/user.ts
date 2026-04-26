/**
 * 用户管理相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 用户信息接口（对应后端 UserResp）
 * 后端返回 mfa_enabled（非 locked）
 */
export interface User {
  id: number
  email: string
  nickname: string
  role: string
  language: string
  timezone: string
  mfa_enabled: boolean
  created_at: string
  updated_at: string
}

/**
 * 更新用户请求接口（对应后端 AdminUpdateUserReq）
 */
export interface UpdateUserReq {
  nickname?: string
  role?: string
  language?: string
  timezone?: string
}

/**
 * 用户列表查询请求接口
 */
export interface UserListReq {
  search?: string
  page?: number
  page_size?: number
}
