/**
 * 用户管理API接口
 * 提供管理员对系统用户的查询、修改、删除、角色变更和锁定/解锁功能
 * 这些接口通常需要管理员权限才能调用
 */
import { get, put, del, post } from '@/utils/request'
import type { User, UpdateUserReq, UserListReq } from '@/types/user'
import type { PageResult } from '@/types/common'

/**
 * 获取用户列表（分页）
 * 支持按关键词搜索用户
 * @param params - 查询参数（可选：搜索关键词、页码、每页数量）
 * @returns 用户分页列表
 * @endpoint GET /users
 */
export function list(params: UserListReq) {
  return get<PageResult<User>>('/users', params as Record<string, unknown>)
}

/**
 * 获取单个用户详情
 * @param id - 用户ID
 * @returns 用户详细信息
 * @endpoint GET /users/:id
 */
export function getUser(id: string) { return get<User>(`/users/${id}`) }

/**
 * 更新用户信息
 * @param id - 用户ID
 * @param data - 更新内容（昵称、角色、语言、时区，均为可选）
 * @returns 更新后的用户信息
 * @endpoint PUT /users/:id
 */
export function update(id: string, data: UpdateUserReq) { return put<User>(`/users/${id}`, data) }

/**
 * 删除用户
 * @param id - 用户ID
 * @endpoint DELETE /users/:id
 */
export function remove(id: string) { return del<void>(`/users/${id}`) }

/**
 * 变更用户角色
 * @param id - 用户ID
 * @param role - 新角色名称（如 'admin'、'user' 等）
 * @endpoint PUT /users/:id/role
 */
export function changeRole(id: string, role: string) { return put<void>(`/users/${id}/role`, { role }) }

/**
 * 锁定用户账户
 * 被锁定的用户将无法登录系统
 * @param id - 用户ID
 * @endpoint POST /users/:id/lock
 */
export function lock(id: string) { return post<void>(`/users/${id}/lock`) }

/**
 * 解锁用户账户
 * 恢复被锁定用户的登录权限
 * @param id - 用户ID
 * @endpoint POST /users/:id/unlock
 */
export function unlock(id: string) { return post<void>(`/users/${id}/unlock`) }
