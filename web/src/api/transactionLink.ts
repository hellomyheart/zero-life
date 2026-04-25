/**
 * 交易关联 API 接口
 * 提供交易关联的查询、创建和删除功能
 */
import { get, post, del } from '@/utils/request'
import type { TransactionLink, CreateTransactionLinkRequest, TransactionLinkListParams } from '@/types/transactionLink'
import type { PageResult } from '@/types/common'

/**
 * 获取交易关联列表（分页）
 * 支持按交易ID筛选
 * @param params - 查询参数（可选：交易ID、页码、每页数量）
 * @returns 交易关联分页列表
 * @endpoint GET /transaction-links
 */
export function list(params?: TransactionLinkListParams) {
  return get<PageResult<TransactionLink>>('/transaction-links', params as Record<string, unknown>)
}

/**
 * 创建交易关联
 * @param data - 创建请求参数（源交易ID、关联类型、目标交易日志ID）
 * @returns 创建的交易关联信息
 * @endpoint POST /transaction-links
 */
export function create(data: CreateTransactionLinkRequest) {
  return post<TransactionLink>('/transaction-links', data)
}

/**
 * 删除交易关联
 * @param id - 关联记录ID
 * @endpoint DELETE /transaction-links/:id
 */
export function remove(id: number) {
  return del<void>(`/transaction-links/${id}`)
}
