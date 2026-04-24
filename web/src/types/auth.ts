/**
 * 认证相关类型定义
 * 包含登录、注册、Token刷新、密码重置和个人信息相关的请求和响应类型
 */

/**
 * 登录请求接口
 * 用户使用邮箱和密码登录系统
 */
export interface LoginReq {
  email: string                // 注册邮箱
  password: string             // 登录密码
}

/**
 * 注册请求接口
 * 新用户注册时需要提供的信息
 */
export interface RegisterReq {
  email: string                // 注册邮箱
  password: string             // 设置密码
  nickname: string             // 用户昵称
}

/**
 * 登录响应接口
 * 登录或注册成功后返回的令牌信息
 */
export interface LoginResp {
  access_token: string         // 访问令牌，用于API请求的身份验证（有效期较短）
  refresh_token: string        // 刷新令牌，用于获取新的访问令牌（有效期较长）
}

/**
 * 用户资料响应接口
 * 获取当前登录用户的个人信息
 */
export interface ProfileResp {
  id: string                   // 用户唯一标识
  email: string                // 用户邮箱
  nickname: string             // 用户昵称
  language: string             // 用户偏好语言，如"zh-CN"、"en-US"
  timezone: string             // 用户所在时区，如"Asia/Shanghai"
  created_at: string           // 账户创建时间
  updated_at: string           // 最后更新时间
}

/**
 * Token刷新请求接口
 * 当访问令牌过期时，使用刷新令牌获取新的令牌对
 */
export interface RefreshReq {
  refresh_token: string        // 刷新令牌
}

/**
 * 忘记密码请求接口
 * 提交注册邮箱，系统会发送密码重置邮件
 */
export interface ForgotPasswordReq {
  email: string                // 注册邮箱地址
}

/**
 * 重置密码请求接口
 * 使用邮件中收到的重置Token来设置新密码
 */
export interface ResetPasswordReq {
  token: string                // 密码重置令牌（从邮件链接中获取）
  new_password: string         // 新密码
}

/**
 * 更新个人资料请求接口
 * 所有字段均为可选，只更新需要修改的字段
 */
export interface UpdateProfileReq {
  nickname?: string            // 新昵称
  language?: string            // 新语言偏好
  timezone?: string            // 新时区设置
}

/**
 * 修改密码请求接口
 * 已登录用户主动修改密码时使用
 */
export interface ChangePasswordReq {
  old_password: string         // 当前密码（用于验证身份）
  new_password: string         // 新密码
}
