import { get, post, put, del } from '@/utils/request'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq, PiggyBankListReq } from '@/types/piggyBank'
import type { PiggyEvent } from '@/types/piggyBank'
import type { PageResult } from '@/types/common'

export function list(params: PiggyBankListReq) {
  return get<PageResult<PiggyBank>>('/piggy-banks', params as Record<string, unknown>)
}
export function getPiggyBank(id: string) { return get<PiggyBank>(`/piggy-banks/${id}`) }
export function create(data: CreatePiggyBankReq) { return post<PiggyBank>('/piggy-banks', data) }
export function update(id: string, data: UpdatePiggyBankReq) { return put<PiggyBank>(`/piggy-banks/${id}`, data) }
export function remove(id: string) { return del<void>(`/piggy-banks/${id}`) }
export function addAmount(id: string, data: AddAmountReq) { return post<void>(`/piggy-banks/${id}/add`, data) }
export function removeAmount(id: string, data: AddAmountReq) { return post<void>(`/piggy-banks/${id}/remove`, data) }
export function getEvents(id: string) { return get<PiggyEvent[]>(`/piggy-banks/${id}/events`) }
