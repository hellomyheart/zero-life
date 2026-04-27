/**
 * 认证相关API接口
 * 提供用户注册、登录、Token刷新、密码重置、个人信息查看与修改等功能
 */
import { post, get, put } from '@/utils/request'
import type { LoginReq, RegisterReq, RefreshReq, ForgotPasswordReq, ResetPasswordReq, UpdateProfileReq, ChangePasswordReq, LoginResp, MFALoginVerifyReq, ProfileResp } from '@/types/auth'

/**
 * 用户注册
 * 提交邮箱、密码和昵称，注册成功后自动登录并返回Token
 * @param data - 注册请求参数（邮箱、密码、昵称）
 * @returns 登录响应（包含 access_token 和 refresh_token）
 * @endpoint POST /auth/register
 */
export function register(data: RegisterReq) {
  return post<LoginResp>('/auth/register', data)
}

/**
 * 用户登录
 * 使用邮箱和密码登录，成功后返回访问令牌和刷新令牌
 * @param data - 登录请求参数（邮箱、密码）
 * @returns 登录响应（包含 access_token 和 refresh_token）
 * @endpoint POST /auth/login
 */
export function login(data: LoginReq) {
  return post<LoginResp>('/auth/login', data)
}

/**
 * 刷新访问令牌
 * 当 access_token 过期时，使用 refresh_token 获取新的令牌对
 * @param data - 刷新请求参数（包含 refresh_token）
 * @returns 新的登录响应（包含新的 access_token 和 refresh_token）
 * @endpoint POST /auth/refresh
 */
export function refresh(data: RefreshReq) {
  return post<LoginResp>('/auth/refresh', data)
}

/**
 * 忘记密码
 * 提交注册邮箱，后端会发送密码重置邮件
 * @param data - 忘记密码请求参数（邮箱地址）
 * @endpoint POST /auth/forgot-password
 */
export function forgotPassword(data: ForgotPasswordReq) {
  return post<void>('/auth/forgot-password', data)
}

/**
 * 重置密码
 * 使用邮件中收到的重置Token设置新密码
 * @param data - 重置密码请求参数（重置Token + 新密码）
 * @endpoint POST /auth/reset-password
 */
export function resetPassword(data: ResetPasswordReq) {
  return post<void>('/auth/reset-password', data)
}

/**
 * 获取当前登录用户的个人信息
 * @returns 用户个人资料（ID、邮箱、昵称、语言、时区等）
 * @endpoint GET /auth/profile
 */
export function getProfile() {
  return get<ProfileResp>('/auth/profile')
}

/**
 * 更新当前登录用户的个人信息
 * @param data - 更新资料请求参数（昵称、语言、时区等，均为可选）
 * @returns 更新后的用户个人资料
 * @endpoint PUT /auth/profile
 */
export function updateProfile(data: UpdateProfileReq) {
  return put<ProfileResp>('/auth/profile', data)
}

/**
 * 修改密码
 * 需要提供旧密码和新密码，用于已登录用户主动修改密码
 * @param data - 修改密码请求参数（旧密码 + 新密码）
 * @endpoint PUT /auth/change-password
 */
export function changePassword(data: ChangePasswordReq) {
  return put<void>('/auth/password', data)
}

/**
 * MFA登录二次验证
 * 使用登录时返回的临时令牌和TOTP验证码，换取真正的访问令牌
 * @param data - MFA验证请求参数（临时令牌 + TOTP验证码）
 * @returns 登录响应（包含 access_token 和 refresh_token）
 * @endpoint POST /auth/mfa-verify
 */
export function mfaLoginVerify(data: MFALoginVerifyReq) {
  return post<LoginResp>('/auth/mfa-verify', data)
}
