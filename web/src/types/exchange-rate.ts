/**
 * 汇率相关类型定义
 * 用于汇率管理模块，包含汇率和货币的详细信息
 */

/**
 * 汇率信息接口
 * 表示两种货币之间在特定日期的汇率
 */
export interface ExchangeRate {
  id: number                   // 汇率记录唯一标识
  user_id: number              // 所属用户ID
  from_currency_id: number     // 源货币ID
  to_currency_id: number       // 目标货币ID
  date: string                 // 汇率日期
  rate: string                 // 系统汇率（自动获取或默认值）
  user_rate?: string           // 用户自定义汇率（优先使用此值）
  from_currency?: Currency     // 源货币详细信息（关联查询时返回）
  to_currency?: Currency       // 目标货币详细信息（关联查询时返回）
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 货币信息接口
 * 汇率模块中使用的货币基本信息
 */
export interface Currency {
  id: number                   // 货币唯一标识
  code: string                 // 货币代码，如"CNY"、"USD"
  name: string                 // 货币名称，如"人民币"
  symbol: string               // 货币符号，如"¥"
  decimal_places: number       // 小数位数（如CNY为2，JPY为0）
}
