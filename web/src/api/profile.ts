/**
 * 用户资料 API 接口
 * 提供当前登录用户个人资料的查看、修改和密码变更功能
 * 所有接口路径与后端 /api/v1/auth/profile 对齐
 */
import { get, put } from '@/utils/request'
import type { ProfileResp, UpdateProfileReq, ChangePasswordReq } from '@/types/auth'

/**
 * 获取当前登录用户的个人资料
 * @returns 用户资料信息
 * @endpoint GET /auth/profile
 */
export function getProfile() {
  return get<ProfileResp>('/auth/profile')
}

/**
 * 更新当前登录用户的个人资料
 * @param data - 更新内容（昵称、语言、时区等，均为可选）
 * @returns 更新后的用户资料信息
 * @endpoint PUT /auth/profile
 */
export function updateProfile(data: UpdateProfileReq) {
  return put<ProfileResp>('/auth/profile', data)
}

/**
 * 修改当前登录用户的密码
 * @param data - 密码修改信息（旧密码 + 新密码）
 * @endpoint PUT /auth/password
 */
export function changePassword(data: ChangePasswordReq) {
  return put<void>('/auth/password', data)
}
