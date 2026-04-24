import { get, put, del, post } from '@/utils/request'
import type { User, UpdateUserReq, UserListReq } from '@/types/user'
import type { PageResult } from '@/types/common'

export function list(params: UserListReq) {
  return get<PageResult<User>>('/users', params as Record<string, unknown>)
}
export function getUser(id: string) { return get<User>(`/users/${id}`) }
export function update(id: string, data: UpdateUserReq) { return put<User>(`/users/${id}`, data) }
export function remove(id: string) { return del<void>(`/users/${id}`) }
export function changeRole(id: string, role: string) { return put<void>(`/users/${id}/role`, { role }) }
export function lock(id: string) { return post<void>(`/users/${id}/lock`) }
export function unlock(id: string) { return post<void>(`/users/${id}/unlock`) }
