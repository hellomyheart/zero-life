import { get, post, put, del } from '@/utils/request'
import type { Transaction, CreateTransactionReq, UpdateTransactionReq, TransactionListReq } from '@/types/transaction'
import type { PageResult } from '@/types/common'

export function list(params: TransactionListReq) {
  return get<PageResult<Transaction>>('/transactions', params as Record<string, unknown>)
}

export function getTransaction(id: string) {
  return get<Transaction>(`/transactions/${id}`)
}

export function create(data: CreateTransactionReq) {
  return post<Transaction>('/transactions', data)
}

export function update(id: string, data: UpdateTransactionReq) {
  return put<Transaction>(`/transactions/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/transactions/${id}`)
}

export function search(params: TransactionListReq) {
  return get<PageResult<Transaction>>('/transactions/search', params as Record<string, unknown>)
}
