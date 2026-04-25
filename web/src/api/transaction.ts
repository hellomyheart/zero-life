/**
 * 交易相关API
 * 提供交易的增删改查、搜索等接口
 */
import { get, post, put, del } from '@/utils/request'
import type { Transaction, CreateTransactionReq, UpdateTransactionReq, TransactionListReq } from '@/types/transaction'
import type { PageResult } from '@/types/common'

/**
 * 获取交易列表
 * @param params 查询参数（类型、日期范围、账户、分类等）
 * @returns 交易分页列表
 */
export function list(params: TransactionListReq) {
  return get<PageResult<Transaction>>('/transactions', params as Record<string, unknown>)
}

/**
 * 获取单个交易详情
 * @param id 交易ID
 * @returns 交易详情
 */
export function getTransaction(id: number) {
  return get<Transaction>(`/transactions/${id}`)
}

/**
 * 创建新交易
 * @param data 交易数据
 * @returns 创建的交易信息
 */
export function create(data: CreateTransactionReq) {
  return post<Transaction>('/transactions', data)
}

/**
 * 更新交易
 * @param id 交易ID
 * @param data 更新数据
 * @returns 更新后的交易信息
 */
export function update(id: number, data: UpdateTransactionReq) {
  return put<Transaction>(`/transactions/${id}`, data)
}

/**
 * 删除交易
 * @param id 交易ID
 */
export function remove(id: number) {
  return del<void>(`/transactions/${id}`)
}

/**
 * 搜索交易
 * @param params 搜索参数
 * @returns 匹配的交易分页列表
 */
export function search(params: TransactionListReq) {
  return get<PageResult<Transaction>>('/transactions/search', params as Record<string, unknown>)
}
