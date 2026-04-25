/**
 * 汇率管理 API 接口
 * 提供汇率的增删改查和货币转换功能
 */
import { get, post, put, del } from '@/utils/request'
import type { ExchangeRate, ExchangeRateListParams, CreateExchangeRateReq, UpdateExchangeRateReq, ConvertResult } from '@/types/exchange-rate'
import type { PageResult } from '@/types/common'

/**
 * 获取汇率列表（分页）
 * 支持按货币对、日期范围筛选
 * @param params - 查询参数（可选：货币ID、日期范围、页码、每页数量）
 * @returns 汇率分页列表
 * @endpoint GET /exchange-rates
 */
export function list(params: ExchangeRateListParams) {
  return get<PageResult<ExchangeRate>>('/exchange-rates', params as Record<string, unknown>)
}

/**
 * 获取汇率详情
 * @param id - 汇率记录ID
 * @returns 汇率详细信息
 * @endpoint GET /exchange-rates/:id
 */
export function getRate(id: number) {
  return get<ExchangeRate>(`/exchange-rates/${id}`)
}

/**
 * 创建汇率
 * @param data - 创建汇率请求参数（源货币、目标货币、日期、汇率值）
 * @returns 创建成功的汇率信息
 * @endpoint POST /exchange-rates
 */
export function create(data: CreateExchangeRateReq) {
  return post<ExchangeRate>('/exchange-rates', data)
}

/**
 * 更新汇率
 * @param id - 要更新的汇率记录ID
 * @param data - 更新内容（源货币、目标货币、日期、汇率值，均为可选）
 * @returns 更新后的汇率信息
 * @endpoint PUT /exchange-rates/:id
 */
export function update(id: number, data: UpdateExchangeRateReq) {
  return put<ExchangeRate>(`/exchange-rates/${id}`, data)
}

/**
 * 删除汇率
 * @param id - 要删除的汇率记录ID
 * @endpoint DELETE /exchange-rates/:id
 */
export function remove(id: number) {
  return del<void>(`/exchange-rates/${id}`)
}

/**
 * 获取最新汇率
 * @param from - 源货币代码
 * @param to - 目标货币代码
 * @returns 最新汇率信息
 * @endpoint GET /exchange-rates/latest/:from/:to
 */
export function getLatest(from: string, to: string) {
  return get<ExchangeRate>(`/exchange-rates/latest/${from}/${to}`)
}

/**
 * 货币转换
 * @param from - 源货币代码
 * @param to - 目标货币代码
 * @param amount - 转换金额
 * @param date - 可选，指定日期的汇率（默认使用最新汇率）
 * @returns 转换结果（包含转换后金额和使用的汇率）
 * @endpoint GET /exchange-rates/convert/:from/:to
 */
export function convert(from: string, to: string, amount: string, date?: string) {
  const params: Record<string, unknown> = { amount }
  if (date) params.date = date
  return get<ConvertResult>(`/exchange-rates/convert/${from}/${to}`, params)
}
