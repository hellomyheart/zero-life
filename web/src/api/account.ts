/**
 * 账户相关API接口
 * 提供账户的增删改查功能，账户是记账系统的基础实体，
 * 包括资产账户（如银行卡）、支出账户、收入账户和负债账户
 */
import { get, post, put, del } from '@/utils/request'
import type { Account, CreateAccountReq, UpdateAccountReq, AccountListReq } from '@/types/account'
import type { PageResult } from '@/types/common'

/**
 * 获取账户列表（分页）
 * 支持按账户类型和货币进行筛选
 * @param params - 查询参数（可选：账户类型、货币、页码、每页数量）
 * @returns 账户分页列表
 * @endpoint GET /accounts
 */
export function list(params: AccountListReq) {
  return get<PageResult<Account>>('/accounts', params as Record<string, unknown>)
}

/**
 * 获取单个账户详情
 * @param id - 账户ID
 * @returns 账户详细信息
 * @endpoint GET /accounts/:id
 */
export function getAccount(id: string) {
  return get<Account>(`/accounts/${id}`)
}

/**
 * 创建新账户
 * @param data - 创建账户请求参数（名称、类型、货币、初始余额等）
 * @returns 创建成功的账户信息
 * @endpoint POST /accounts
 */
export function create(data: CreateAccountReq) {
  return post<Account>('/accounts', data)
}

/**
 * 更新账户信息
 * @param id - 要更新的账户ID
 * @param data - 更新内容（名称、类型、货币、初始余额等，均为可选）
 * @returns 更新后的账户信息
 * @endpoint PUT /accounts/:id
 */
export function update(id: string, data: UpdateAccountReq) {
  return put<Account>(`/accounts/${id}`, data)
}

/**
 * 删除账户
 * @param id - 要删除的账户ID
 * @endpoint DELETE /accounts/:id
 */
export function remove(id: string) {
  return del<void>(`/accounts/${id}`)
}
