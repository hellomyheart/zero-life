import { get, post, put, del } from '@/utils/request'
import type { Bill, CreateBillReq, UpdateBillReq } from '@/types/bill'

export function list() {
  return get<Bill[]>('/bills')
}

export function getBill(id: string) {
  return get<Bill>(`/bills/${id}`)
}

export function create(data: CreateBillReq) {
  return post<Bill>('/bills', data)
}

export function update(id: string, data: UpdateBillReq) {
  return put<Bill>(`/bills/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/bills/${id}`)
}
