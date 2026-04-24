/**
 * 汇率相关类型定义
 */

/**
 * 汇率信息
 */
export interface ExchangeRate {
  id: number
  user_id: number
  from_currency_id: number
  to_currency_id: number
  date: string
  rate: string
  user_rate?: string
  from_currency?: Currency
  to_currency?: Currency
  created_at: string
  updated_at: string
}

/**
 * 货币信息
 */
export interface Currency {
  id: number
  code: string
  name: string
  symbol: string
  decimal_places: number
}
