/**
 * 用户偏好 API 接口
 * 提供用户偏好的查询、设置和删除功能
 */
import { get, put, del } from '@/utils/request'
import type { Preference, SetPreferenceRequest } from '@/types/preference'

/**
 * 获取所有用户偏好
 * @returns 偏好列表
 * @endpoint GET /preferences
 */
export function list() {
  return get<Preference[]>('/preferences')
}

/**
 * 获取单个偏好值
 * @param key - 偏好键名（路径参数）
 * @returns 偏好信息
 * @endpoint GET /preferences/:key
 */
export function getPreference(key: string) {
  return get<Preference>(`/preferences/${key}`)
}

/**
 * 设置偏好值（创建或更新）
 * @param data - 设置请求参数（键名和值）
 * @returns 设置后的偏好信息
 * @endpoint PUT /preferences
 */
export function set(data: SetPreferenceRequest) {
  return put<Preference>('/preferences', data)
}

/**
 * 删除偏好
 * @param key - 偏好键名
 * @endpoint DELETE /preferences/:key
 */
export function remove(key: string) {
  return del<void>(`/preferences/${key}`)
}
