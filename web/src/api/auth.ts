// 认证相关API接口 - 登录、注册、Token刷新、密码重置、个人信息
import { post, get, put } from '@/utils/request'
import type { LoginReq, RegisterReq, RefreshReq, ForgotPasswordReq, ResetPasswordReq, UpdateProfileReq, ChangePasswordReq, LoginResp, ProfileResp } from '@/types/auth'

export function register(data: RegisterReq) {
  return post<LoginResp>('/auth/register', data)
}

export function login(data: LoginReq) {
  return post<LoginResp>('/auth/login', data)
}

export function refresh(data: RefreshReq) {
  return post<LoginResp>('/auth/refresh', data)
}

export function forgotPassword(data: ForgotPasswordReq) {
  return post<void>('/auth/forgot-password', data)
}

export function resetPassword(data: ResetPasswordReq) {
  return post<void>('/auth/reset-password', data)
}

export function getProfile() {
  return get<ProfileResp>('/auth/profile')
}

export function updateProfile(data: UpdateProfileReq) {
  return put<ProfileResp>('/auth/profile', data)
}

export function changePassword(data: ChangePasswordReq) {
  return put<void>('/auth/change-password', data)
}
