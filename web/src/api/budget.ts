/**
 * 预算相关API接口
 * 提供预算的增删改查和历史记录查询功能
 * 预算用于设定某个时间段内（月/季/年）特定分类的支出上限，帮助控制开支
 */
import { get, post, put, del } from '@/utils/request'
import type { Budget, CreateBudgetReq, UpdateBudgetReq, BudgetHistory } from '@/types/budget'

/**
 * 获取所有预算列表
 * @returns 预算数组
 * @endpoint GET /budgets
 */
export function list() {
  return get<Budget[]>('/budgets')
}

/**
 * 获取单个预算详情
 * @param id - 预算ID
 * @returns 预算详细信息（含已花费金额和使用率）
 * @endpoint GET /budgets/:id
 */
export function getBudget(id: string) {
  return get<Budget>(`/budgets/${id}`)
}

/**
 * 创建新预算
 * @param data - 创建预算请求参数（名称、金额、周期、关联分类、开始日期）
 * @returns 创建成功的预算信息
 * @endpoint POST /budgets
 */
export function create(data: CreateBudgetReq) {
  return post<Budget>('/budgets', data)
}

/**
 * 更新预算信息
 * @param id - 要更新的预算ID
 * @param data - 更新内容（名称、金额、周期、关联分类、开始日期，均为可选）
 * @returns 更新后的预算信息
 * @endpoint PUT /budgets/:id
 */
export function update(id: string, data: UpdateBudgetReq) {
  return put<Budget>(`/budgets/${id}`, data)
}

/**
 * 删除预算
 * @param id - 要删除的预算ID
 * @endpoint DELETE /budgets/:id
 */
export function remove(id: string) {
  return del<void>(`/budgets/${id}`)
}

/**
 * 获取预算的历史使用记录
 * 返回该预算在各个周期内的花费情况和使用率
 * @param id - 预算ID
 * @returns 预算历史记录数组
 * @endpoint GET /budgets/:id/history
 */
export function getHistory(id: string) {
  return get<BudgetHistory[]>(`/budgets/${id}/history`)
}
