/**
 * 循环交易相关API接口
 * 提供循环交易的增删改查和到期处理功能
 * 循环交易用于记录定期重复发生的交易（如每月工资、每周通勤费等），
 * 系统会在到期时自动创建对应的交易记录
 */
import { get, post, put, del } from '@/utils/request'
import type { RecurringTransaction, CreateRecurringTransactionReq, UpdateRecurringTransactionReq, RecurringTransactionListReq } from '@/types/recurringTransaction'
import type { PageResult } from '@/types/common'

/**
 * 获取循环交易列表（分页）
 * @param params - 查询参数（可选：是否仅显示活跃的、页码、每页数量）
 * @returns 循环交易分页列表
 * @endpoint GET /recurring-transactions
 */
export function list(params: RecurringTransactionListReq) {
  return get<PageResult<RecurringTransaction>>('/recurring-transactions', params as Record<string, unknown>)
}

/**
 * 获取单个循环交易详情
 * @param id - 循环交易ID
 * @returns 循环交易详细信息
 * @endpoint GET /recurring-transactions/:id
 */
export function getRecurringTransaction(id: string) { return get<RecurringTransaction>(`/recurring-transactions/${id}`) }

/**
 * 创建新循环交易
 * @param data - 创建请求参数（标题、类型、金额、源账户、重复类型、间隔、开始日期等）
 * @returns 创建成功的循环交易信息
 * @endpoint POST /recurring-transactions
 */
export function create(data: CreateRecurringTransactionReq) { return post<RecurringTransaction>('/recurring-transactions', data) }

/**
 * 更新循环交易信息
 * @param id - 循环交易ID
 * @param data - 更新内容（标题、金额、重复类型、间隔、日期、是否活跃等，均为可选）
 * @returns 更新后的循环交易信息
 * @endpoint PUT /recurring-transactions/:id
 */
export function update(id: string, data: UpdateRecurringTransactionReq) { return put<RecurringTransaction>(`/recurring-transactions/${id}`, data) }

/**
 * 删除循环交易
 * @param id - 循环交易ID
 * @endpoint DELETE /recurring-transactions/:id
 */
export function remove(id: string) { return del<void>(`/recurring-transactions/${id}`) }

/**
 * 处理到期的循环交易
 * 系统会检查所有活跃的循环交易，将到期的交易自动创建为正式交易记录
 * @endpoint POST /recurring-transactions/process-due
 */
export function processDue() { return post<void>('/recurring-transactions/process-due') }
