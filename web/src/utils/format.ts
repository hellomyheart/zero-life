import dayjs from 'dayjs'
import type { Currency } from '@/types/currency'

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

export function formatDate(date: string | Date, format: string = 'YYYY-MM-DD'): string {
  return dayjs(date).format(format)
}

export function formatPercent(rate: number): string {
  return (rate * 100).toFixed(2) + '%'
}
