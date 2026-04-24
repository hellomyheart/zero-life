/**
 * 用户偏好 API 接口
 * 提供用户偏好的 CRUD 操作
 */
import request from '@/utils/request';
import type {
  Preference,
  SetPreferenceRequest,
  GetPreferenceParams
} from '@/types/preference';

/**
 * 获取所有用户偏好
 * @returns 偏好列表
 */
export function list() {
  return request<Preference[]>({
    url: '/preferences',
    method: 'get'
  });
}

/**
 * 获取单个偏好值
 * @param params 查询参数（包含 key）
 * @returns 偏好值
 */
export function get(params: GetPreferenceParams) {
  return request<Preference>({
    url: '/preferences',
    method: 'get',
    params
  });
}

/**
 * 设置偏好值
 * @param data 设置请求参数
 * @returns 设置后的偏好
 */
export function set(data: SetPreferenceRequest) {
  return request<Preference>({
    url: '/preferences',
    method: 'post',
    data
  });
}

/**
 * 删除偏好
 * @param key 偏好键名
 */
export function remove(key: string) {
  return request({
    url: `/preferences/${key}`,
    method: 'delete'
  });
}
