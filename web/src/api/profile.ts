/**
 * 用户资料 API 接口
 */
import request from '@/utils/request'

/**
 * 获取用户资料
 */
export function getProfile() {
  return request.get('/profile')
}

/**
 * 更新用户资料
 */
export function updateProfile(data: any) {
  return request.put('/profile', data)
}

/**
 * 修改密码
 */
export function changePassword(data: any) {
  return request.post('/profile/password', data)
}
