// 循环交易API接口 - 循环交易增删改查和到期提醒
import { get, post, put, del } from '@/utils/request'
import type { RecurringTransaction, CreateRecurringTransactionReq, UpdateRecurringTransactionReq, RecurringTransactionListReq } from '@/types/recurringTransaction'
import type { PageResult } from '@/types/common'

export function list(params: RecurringTransactionListReq) {
  return get<PageResult<RecurringTransaction>>('/recurring-transactions', params as Record<string, unknown>)
}
export function getRecurringTransaction(id: string) { return get<RecurringTransaction>(`/recurring-transactions/${id}`) }
export function create(data: CreateRecurringTransactionReq) { return post<RecurringTransaction>('/recurring-transactions', data) }
export function update(id: string, data: UpdateRecurringTransactionReq) { return put<RecurringTransaction>(`/recurring-transactions/${id}`, data) }
export function remove(id: string) { return del<void>(`/recurring-transactions/${id}`) }
export function processDue() { return post<void>('/recurring-transactions/process-due') }
