/**
 * 认证相关类型定义
 * 包含登录、注册、Token刷新、密码重置和个人信息相关的请求和响应类型
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 登录请求接口（对应后端 LoginReq）
 */
export interface LoginReq {
  email: string
  password: string
}

/**
 * 注册请求接口（对应后端 RegisterReq）
 */
export interface RegisterReq {
  email: string
  password: string
  nickname: string
}

/**
 * 登录响应接口（对应后端 LoginResp）
 */
export interface LoginResp {
  access_token: string
  refresh_token: string
  expires_at: string
}

/**
 * 用户资料响应接口（对应后端 ProfileResp）
 * 后端不返回 updated_at，仅返回 id/email/nickname/language/timezone/created_at
 */
export interface ProfileResp {
  id: number
  email: string
  nickname: string
  language: string
  timezone: string
  created_at: string
}

/**
 * Token刷新请求接口（对应后端 RefreshReq）
 */
export interface RefreshReq {
  refresh_token: string
}

/**
 * 忘记密码请求接口（对应后端 ForgotPasswordReq）
 */
export interface ForgotPasswordReq {
  email: string
}

/**
 * 重置密码请求接口（对应后端 ResetPasswordReq）
 * 后端字段名为 password（非 new_password）
 */
export interface ResetPasswordReq {
  token: string
  password: string
}

/**
 * 更新个人资料请求接口（对应后端 UpdateProfileReq）
 * 后端仅接受 nickname/language/timezone，不接受 email
 */
export interface UpdateProfileReq {
  nickname?: string
  language?: string
  timezone?: string
}

/**
 * 修改密码请求接口（对应后端 ChangePasswordReq）
 */
export interface ChangePasswordReq {
  old_password: string
  new_password: string
}
