/**
 * 分类相关API接口
 * 提供分类的增删改查功能，分类用于对交易进行归类（如餐饮、交通、工资等）
 * 分类支持树形结构，即可以有父子级关系
 */
import { get, post, put, del } from '@/utils/request'
import type { Category, CreateCategoryReq, UpdateCategoryReq } from '@/types/category'

/**
 * 获取所有分类列表
 * 返回完整的分类树形结构，包含子分类
 * @returns 分类数组（含嵌套的子分类）
 * @endpoint GET /categories
 */
export function list() {
  return get<Category[]>('/categories')
}

/**
 * 获取单个分类详情
 * @param id - 分类ID
 * @returns 分类详细信息
 * @endpoint GET /categories/:id
 */
export function getCategory(id: number) {
  return get<Category>(`/categories/${id}`)
}

/**
 * 创建新分类
 * @param data - 创建分类请求参数（名称、可选的父分类ID）
 * @returns 创建成功的分类信息
 * @endpoint POST /categories
 */
export function create(data: CreateCategoryReq) {
  return post<Category>('/categories', data)
}

/**
 * 更新分类信息
 * @param id - 要更新的分类ID
 * @param data - 更新内容（名称、父分类ID，均为可选）
 * @returns 更新后的分类信息
 * @endpoint PUT /categories/:id
 */
export function update(id: number, data: UpdateCategoryReq) {
  return put<Category>(`/categories/${id}`, data)
}

/**
 * 删除分类
 * @param id - 要删除的分类ID
 * @endpoint DELETE /categories/:id
 */
export function remove(id: number) {
  return del<void>(`/categories/${id}`)
}
