/**
 * 货币相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 货币信息接口（对应后端 CurrencyResp）
 */
export interface Currency {
  id: number
  code: string
  name: string
  symbol: string
  decimal_places: number
  is_enabled: boolean
  is_default: boolean
  created_at: string
  updated_at: string
}

/**
 * 汇率信息接口（对应后端 ExchangeRateResp）
 * 后端返回 from_currency_id/to_currency_id（非 source_currency/target_currency）
 */
export interface ExchangeRate {
  id: number
  from_currency_id: number
  to_currency_id: number
  rate: string
  updated_at: string
}

/**
 * 设置货币启用状态请求接口
 */
export interface SetCurrencyStatusReq {
  is_enabled: boolean
}

/**
 * 设置汇率请求接口（对应后端 SetExchangeRateReq）
 */
export interface SetExchangeRateReq {
  from_currency_id: number
  to_currency_id: number
  rate: string
}
