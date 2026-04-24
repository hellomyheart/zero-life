/**
 * 对账相关API接口
 * 提供对账记录的增删改查功能
 * 对账用于核对账户的账面余额与实际余额是否一致，帮助发现记账错误或遗漏
 */
import { get, post, put, del } from '@/utils/request'
import type { Reconciliation, CreateReconciliationReq, UpdateReconciliationReq, ReconciliationListReq } from '@/types/reconciliation'
import type { PageResult } from '@/types/common'

/**
 * 获取对账记录列表（分页）
 * @param params - 查询参数（可选：账户ID、页码、每页数量）
 * @returns 对账记录分页列表
 * @endpoint GET /reconciliations
 */
export function list(params: ReconciliationListReq) {
  return get<PageResult<Reconciliation>>('/reconciliations', params as Record<string, unknown>)
}

/**
 * 获取单条对账记录详情
 * @param id - 对账记录ID
 * @returns 对账记录详细信息
 * @endpoint GET /reconciliations/:id
 */
export function getReconciliation(id: string) { return get<Reconciliation>(`/reconciliations/${id}`) }

/**
 * 创建新对账记录
 * @param data - 创建请求参数（账户ID、起止日期、起止余额、备注）
 * @returns 创建成功的对账记录
 * @endpoint POST /reconciliations
 */
export function create(data: CreateReconciliationReq) { return post<Reconciliation>('/reconciliations', data) }

/**
 * 更新对账记录
 * @param id - 对账记录ID
 * @param data - 更新内容（期末余额、状态、备注，均为可选）
 * @returns 更新后的对账记录
 * @endpoint PUT /reconciliations/:id
 */
export function update(id: string, data: UpdateReconciliationReq) { return put<Reconciliation>(`/reconciliations/${id}`, data) }

/**
 * 删除对账记录
 * @param id - 对账记录ID
 * @endpoint DELETE /reconciliations/:id
 */
export function remove(id: string) { return del<void>(`/reconciliations/${id}`) }
