// 标签相关API接口 - 标签的增删改查
import { get, post, put, del } from '@/utils/request'
import type { Tag, CreateTagReq, UpdateTagReq } from '@/types/tag'

export function list() {
  return get<Tag[]>('/tags')
}

export function getTag(id: string) {
  return get<Tag>(`/tags/${id}`)
}

export function create(data: CreateTagReq) {
  return post<Tag>('/tags', data)
}

export function update(id: string, data: UpdateTagReq) {
  return put<Tag>(`/tags/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/tags/${id}`)
}
