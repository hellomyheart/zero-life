/**
 * 汇率管理 API 接口
 * 后端汇率功能集成在 /currencies/exchange-rates 下
 * 仅提供汇率列表查询和设置功能
 */
import { get, post } from '@/utils/request'
import type { ExchangeRate, SetExchangeRateReq } from '@/types/currency'

/**
 * 获取所有汇率列表
 * @returns 汇率数组
 * @endpoint GET /currencies/exchange-rates
 */
export function list() {
  return get<ExchangeRate[]>('/currencies/exchange-rates')
}

/**
 * 设置汇率
 * 创建或更新两种货币之间的汇率
 * @param data - 汇率设置请求（源货币ID、目标货币ID、汇率值）
 * @returns 汇率信息
 * @endpoint POST /currencies/exchange-rates
 */
export function set(data: SetExchangeRateReq) {
  return post<ExchangeRate>('/currencies/exchange-rates', data)
}
