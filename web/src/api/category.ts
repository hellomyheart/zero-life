// 分类相关API接口 - 分类的增删改查
import { get, post, put, del } from '@/utils/request'
import type { Category, CreateCategoryReq, UpdateCategoryReq } from '@/types/category'

export function list() {
  return get<Category[]>('/categories')
}

export function getCategory(id: string) {
  return get<Category>(`/categories/${id}`)
}

export function create(data: CreateCategoryReq) {
  return post<Category>('/categories', data)
}

export function update(id: string, data: UpdateCategoryReq) {
  return put<Category>(`/categories/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/categories/${id}`)
}
