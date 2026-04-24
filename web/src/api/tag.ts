/**
 * 标签相关API接口
 * 提供标签的增删改查功能，标签用于对交易进行标记和分组（如"出差"、"日常"等）
 * 与分类不同，标签是一种更灵活的横向归类方式，一笔交易可以有多个标签
 */
import { get, post, put, del } from '@/utils/request'
import type { Tag, CreateTagReq, UpdateTagReq } from '@/types/tag'

/**
 * 获取所有标签列表
 * @returns 标签数组
 * @endpoint GET /tags
 */
export function list() {
  return get<Tag[]>('/tags')
}

/**
 * 获取单个标签详情
 * @param id - 标签ID
 * @returns 标签详细信息
 * @endpoint GET /tags/:id
 */
export function getTag(id: string) {
  return get<Tag>(`/tags/${id}`)
}

/**
 * 创建新标签
 * @param data - 创建标签请求参数（名称、颜色）
 * @returns 创建成功的标签信息
 * @endpoint POST /tags
 */
export function create(data: CreateTagReq) {
  return post<Tag>('/tags', data)
}

/**
 * 更新标签信息
 * @param id - 要更新的标签ID
 * @param data - 更新内容（名称、颜色，均为可选）
 * @returns 更新后的标签信息
 * @endpoint PUT /tags/:id
 */
export function update(id: string, data: UpdateTagReq) {
  return put<Tag>(`/tags/${id}`, data)
}

/**
 * 删除标签
 * @param id - 要删除的标签ID
 * @endpoint DELETE /tags/:id
 */
export function remove(id: string) {
  return del<void>(`/tags/${id}`)
}
