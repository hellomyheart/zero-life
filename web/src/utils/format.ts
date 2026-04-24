// 格式化工具 - 金额、日期、百分比的格式化函数
import dayjs from 'dayjs'
import type { Currency } from '@/types/currency'

// 货币符号映射表
const currencySymbolMap: Record<string, string> = {
  CNY: '¥',
  USD: '$',
  EUR: '€',
  GBP: '£',
  JPY: '¥',
  KRW: '₩',
  HKD: 'HK$',
  TWD: 'NT$',
  SGD: 'S$',
  AUD: 'A$',
  CAD: 'C$',
}

// 格式化金额 - 添加货币符号和千分位分隔，保留两位小数
export function formatAmount(amount: string | number, currency?: Currency | string): string {
  const num = Number(amount)
  if (isNaN(num)) return String(amount)

  const currencyCode = typeof currency === 'string' ? currency : currency?.code
  const symbol = currencyCode ? (currencySymbolMap[currencyCode] ?? currencyCode + ' ') : ''

  const formatted = Math.abs(num).toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })

  const sign = num < 0 ? '-' : ''
  return `${sign}${symbol}${formatted}`
}

// 格式化日期 - 使用dayjs库按指定格式输出
export function formatDate(date: string | Date, format: string = 'YYYY-MM-DD'): string {
  return dayjs(date).format(format)
}

// 格式化百分比 - 将小数转为百分比字符串，如0.5 → "50.00%"
export function formatPercent(rate: number): string {
  return (rate * 100).toFixed(2) + '%'
}
