/**
 * 对象分组 API 接口
 * 提供对象分组的增删改查功能
 */
import { get, post, put, del } from '@/utils/request'
import type { ObjectGroup, CreateObjectGroupRequest, UpdateObjectGroupRequest, ObjectGroupListParams } from '@/types/objectGroup'

/**
 * 获取对象分组列表
 * 支持按关联对象类型筛选
 * @param params - 查询参数（可选：对象类型、页码、每页数量）
 * @returns 对象分组列表
 * @endpoint GET /object-groups
 */
export function list(params?: ObjectGroupListParams) {
  return get<ObjectGroup[]>('/object-groups', params as Record<string, unknown>)
}

/**
 * 获取单个对象分组详情
 * @param id - 分组ID
 * @returns 对象分组详情
 * @endpoint GET /object-groups/:id
 */
export function get(id: number) {
  return get<ObjectGroup>(`/object-groups/${id}`)
}

/**
 * 创建对象分组
 * @param data - 创建请求参数（分组名称、关联对象类型、关联对象ID）
 * @returns 创建的对象分组信息
 * @endpoint POST /object-groups
 */
export function create(data: CreateObjectGroupRequest) {
  return post<ObjectGroup>('/object-groups', data)
}

/**
 * 更新对象分组
 * @param id - 分组ID
 * @param data - 更新请求参数（分组名称）
 * @returns 更新后的对象分组信息
 * @endpoint PUT /object-groups/:id
 */
export function update(id: number, data: UpdateObjectGroupRequest) {
  return put<ObjectGroup>(`/object-groups/${id}`, data)
}

/**
 * 删除对象分组
 * @param id - 分组ID
 * @endpoint DELETE /object-groups/:id
 */
export function remove(id: number) {
  return del<void>(`/object-groups/${id}`)
}
