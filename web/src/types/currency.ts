/**
 * 货币相关类型定义
 * 包含货币信息和汇率接口，用于支持多币种记账
 */

/**
 * 货币信息接口
 * 表示系统中的一种货币
 */
export interface Currency {
  id: string                   // 货币唯一标识
  code: string                 // 货币代码，如"CNY"、"USD"、"EUR"
  name: string                 // 货币名称，如"人民币"、"美元"
  symbol: string               // 货币符号，如"¥"、"$"、"€"
  is_enabled: boolean          // 是否启用（禁用的货币不会出现在选择列表中）
  is_default: boolean          // 是否为系统默认货币
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 汇率信息接口
 * 表示两种货币之间的汇率
 */
export interface ExchangeRate {
  source_currency: string      // 源货币代码，如"USD"
  target_currency: string      // 目标货币代码，如"CNY"
  rate: string                 // 汇率值，如"7.24"（1 USD = 7.24 CNY）
  updated_at: string           // 汇率更新时间
}

/**
 * 设置货币启用状态请求接口
 */
export interface SetCurrencyStatusReq {
  is_enabled: boolean          // 是否启用该货币
}

/**
 * 设置默认货币请求接口
 */
export interface SetDefaultCurrencyReq {
  currency_id: string          // 要设为默认的货币ID
}

/**
 * 设置汇率请求接口
 * 创建或更新两种货币之间的汇率
 */
export interface SetExchangeRateReq {
  source_currency: string      // 源货币代码
  target_currency: string      // 目标货币代码
  rate: string                 // 汇率值
}
