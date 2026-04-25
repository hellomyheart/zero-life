/**
 * 货币相关API接口
 * 提供货币列表查询、启用/禁用货币、设置默认货币和汇率管理功能
 * 货币是账户和交易的基础，支持多币种记账
 */
import { get, put, post } from '@/utils/request'
import type { Currency, ExchangeRate, SetCurrencyStatusReq, SetDefaultCurrencyReq, SetExchangeRateReq } from '@/types/currency'

/**
 * 获取所有货币列表
 * @returns 货币数组（包含是否启用、是否为默认货币等信息）
 * @endpoint GET /currencies
 */
export function list() {
  return get<Currency[]>('/currencies')
}

/**
 * 更新货币的启用/禁用状态
 * 禁用的货币不会出现在创建账户时的货币选择列表中
 * @param id - 货币ID
 * @param data - 状态设置请求（是否启用）
 * @returns 更新后的货币信息
 * @endpoint PUT /currencies/:id/status
 */
export function updateStatus(id: number, data: SetCurrencyStatusReq) {
  return put<Currency>(`/currencies/${id}/status`, data)
}

/**
 * 设置默认货币
 * 默认货币用于新建账户时的默认选项和系统展示
 * @param data - 默认货币设置请求（货币ID）
 * @returns 更新后的货币信息
 * @endpoint PUT /currencies/default
 */
export function setDefault(data: SetDefaultCurrencyReq) {
  return put<Currency>('/currencies/default', data)
}

/**
 * 获取所有汇率列表
 * @returns 汇率数组（包含源货币、目标货币和汇率值）
 * @endpoint GET /currencies/exchange-rates
 */
export function getExchangeRates() {
  return get<ExchangeRate[]>('/currencies/exchange-rates')
}

/**
 * 设置汇率
 * 创建或更新两种货币之间的汇率
 * @param data - 汇率设置请求（源货币代码、目标货币代码、汇率值）
 * @returns 汇率信息
 * @endpoint POST /currencies/exchange-rates
 */
export function setExchangeRate(data: SetExchangeRateReq) {
  return post<ExchangeRate>('/currencies/exchange-rates', data)
}
