/**
 * 用户管理相关类型定义
 * 用于管理员管理系统中的用户账户
 */

/**
 * 用户信息接口
 * 表示系统中的一个用户
 */
export interface User {
  id: string                   // 用户唯一标识
  email: string                // 用户邮箱（也用作登录名）
  nickname: string             // 用户昵称
  role: string                 // 用户角色，如"admin"（管理员）、"user"（普通用户）
  language: string             // 用户偏好语言
  timezone: string             // 用户所在时区
  locked: boolean              // 账户是否被锁定（锁定的用户无法登录）
  created_at: string           // 账户创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 更新用户请求接口
 * 管理员修改用户信息时使用，所有字段均为可选
 */
export interface UpdateUserReq {
  nickname?: string            // 用户昵称
  role?: string                // 用户角色
  language?: string            // 偏好语言
  timezone?: string            // 所在时区
}

/**
 * 用户列表查询请求接口
 */
export interface UserListReq {
  search?: string              // 搜索关键词（匹配邮箱或昵称）
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}
