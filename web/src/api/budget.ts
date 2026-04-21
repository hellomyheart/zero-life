import { get, post, put, del } from '@/utils/request'
import type { Budget, CreateBudgetReq, UpdateBudgetReq, BudgetHistory } from '@/types/budget'

export function list() {
  return get<Budget[]>('/budgets')
}

export function getBudget(id: string) {
  return get<Budget>(`/budgets/${id}`)
}

export function create(data: CreateBudgetReq) {
  return post<Budget>('/budgets', data)
}

export function update(id: string, data: UpdateBudgetReq) {
  return put<Budget>(`/budgets/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/budgets/${id}`)
}

export function getHistory(id: string) {
  return get<BudgetHistory[]>(`/budgets/${id}/history`)
}
