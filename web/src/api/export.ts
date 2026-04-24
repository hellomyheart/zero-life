/**
 * 数据导出API接口
 * 提供各类数据（交易、账户、预算、分类、标签）的导出功能
 * 导出结果以文件流（blob）形式返回，支持CSV和JSON格式
 */
import { get } from '@/utils/request'
import type { ExportReq } from '@/types/export'

/**
 * 导出交易数据
 * @param params - 导出参数（可选：日期范围、导出格式）
 * @endpoint GET /exports/transactions
 */
export function exportTransactions(params: ExportReq) { return get<void>('/exports/transactions', params as Record<string, unknown>, { responseType: 'blob' }) }

/**
 * 导出账户数据
 * @param params - 导出参数（可选：日期范围、导出格式）
 * @endpoint GET /exports/accounts
 */
export function exportAccounts(params: ExportReq) { return get<void>('/exports/accounts', params as Record<string, unknown>, { responseType: 'blob' }) }

/**
 * 导出预算数据
 * @param params - 导出参数（可选：日期范围、导出格式）
 * @endpoint GET /exports/budgets
 */
export function exportBudgets(params: ExportReq) { return get<void>('/exports/budgets', params as Record<string, unknown>, { responseType: 'blob' }) }

/**
 * 导出分类数据
 * @param params - 导出参数（可选：日期范围、导出格式）
 * @endpoint GET /exports/categories
 */
export function exportCategories(params: ExportReq) { return get<void>('/exports/categories', params as Record<string, unknown>, { responseType: 'blob' }) }

/**
 * 导出标签数据
 * @param params - 导出参数（可选：日期范围、导出格式）
 * @endpoint GET /exports/tags
 */
export function exportTags(params: ExportReq) { return get<void>('/exports/tags', params as Record<string, unknown>, { responseType: 'blob' }) }
