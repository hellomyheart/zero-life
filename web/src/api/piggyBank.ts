import { get, post, put, del } from '@/utils/request'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq, RemoveAmountReq, PiggyEvent } from '@/types/piggyBank'

export function list() {
  return get<PiggyBank[]>('/piggy-banks')
}

export function getPiggyBank(id: number) {
  return get<PiggyBank>(`/piggy-banks/${id}`)
}

export function create(data: CreatePiggyBankReq) {
  return post<PiggyBank>('/piggy-banks', data)
}

export function update(id: number, data: UpdatePiggyBankReq) {
  return put<PiggyBank>(`/piggy-banks/${id}`, data)
}

export function remove(id: number) {
  return del<void>(`/piggy-banks/${id}`)
}

export function addAmount(id: number, data: AddAmountReq) {
  return post<PiggyBank>(`/piggy-banks/${id}/add`, data)
}

export function removeAmount(id: number, data: RemoveAmountReq) {
  return post<PiggyBank>(`/piggy-banks/${id}/remove`, data)
}

export function getEvents(id: number) {
  return get<PiggyEvent[]>(`/piggy-banks/${id}/events`)
}

export function reorder(orders: Record<number, number>) {
  return put<void>('/piggy-banks/reorder', { orders })
}

export function resetHistory(id: number) {
  return post<PiggyBank>(`/piggy-banks/${id}/reset`)
}
