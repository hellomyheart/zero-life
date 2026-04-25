/**
 * 汇率相关类型定义（已废弃，汇率功能合并到 currency.ts）
 * 此文件保留仅为兼容性，新代码请使用 @/types/currency 中的 ExchangeRate 和 SetExchangeRateReq
 */

export type { ExchangeRate, SetExchangeRateReq as CreateExchangeRateReq, SetExchangeRateReq as UpdateExchangeRateReq } from '@/types/currency'