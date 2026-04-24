// 账户相关API接口 - 账户的增删改查
import { get, post, put, del } from '@/utils/request'
import type { Account, CreateAccountReq, UpdateAccountReq, AccountListReq } from '@/types/account'
import type { PageResult } from '@/types/common'

export function list(params: AccountListReq) {
  return get<PageResult<Account>>('/accounts', params as Record<string, unknown>)
}

export function getAccount(id: string) {
  return get<Account>(`/accounts/${id}`)
}

export function create(data: CreateAccountReq) {
  return post<Account>('/accounts', data)
}

export function update(id: string, data: UpdateAccountReq) {
  return put<Account>(`/accounts/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/accounts/${id}`)
}
