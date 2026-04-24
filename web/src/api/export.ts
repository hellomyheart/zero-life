import { get } from '@/utils/request'
import type { ExportReq } from '@/types/export'

export function exportTransactions(params: ExportReq) { return get<void>('/exports/transactions', params as Record<string, unknown>, { responseType: 'blob' }) }
export function exportAccounts(params: ExportReq) { return get<void>('/exports/accounts', params as Record<string, unknown>, { responseType: 'blob' }) }
export function exportBudgets(params: ExportReq) { return get<void>('/exports/budgets', params as Record<string, unknown>, { responseType: 'blob' }) }
export function exportCategories(params: ExportReq) { return get<void>('/exports/categories', params as Record<string, unknown>, { responseType: 'blob' }) }
export function exportTags(params: ExportReq) { return get<void>('/exports/tags', params as Record<string, unknown>, { responseType: 'blob' }) }
