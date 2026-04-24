/**
 * 用户资料API接口
 * 提供当前登录用户个人资料的查看、修改和密码变更功能
 * 注意：此模块与 auth.ts 中的 profile 接口功能重叠，后续可能会合并
 */
import request from '@/utils/request'

/**
 * 获取当前登录用户的个人资料
 * @returns 用户资料信息
 * @endpoint GET /profile
 */
export function getProfile() {
  return request.get('/profile')
}

/**
 * 更新当前登录用户的个人资料
 * @param data - 更新内容（昵称、语言、时区等）
 * @endpoint PUT /profile
 */
export function updateProfile(data: any) {
  return request.put('/profile', data)
}

/**
 * 修改当前登录用户的密码
 * @param data - 密码修改信息（旧密码 + 新密码）
 * @endpoint POST /profile/password
 */
export function changePassword(data: any) {
  return request.post('/profile/password', data)
}
