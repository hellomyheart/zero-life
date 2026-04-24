// 对账API接口 - 对账记录增删改查
import { get, post, put, del } from '@/utils/request'
import type { Reconciliation, CreateReconciliationReq, UpdateReconciliationReq, ReconciliationListReq } from '@/types/reconciliation'
import type { PageResult } from '@/types/common'

export function list(params: ReconciliationListReq) {
  return get<PageResult<Reconciliation>>('/reconciliations', params as Record<string, unknown>)
}
export function getReconciliation(id: string) { return get<Reconciliation>(`/reconciliations/${id}`) }
export function create(data: CreateReconciliationReq) { return post<Reconciliation>('/reconciliations', data) }
export function update(id: string, data: UpdateReconciliationReq) { return put<Reconciliation>(`/reconciliations/${id}`, data) }
export function remove(id: string) { return del<void>(`/reconciliations/${id}`) }
