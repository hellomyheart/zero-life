// 高精度十进制运算工具 - 基于decimal.js封装，避免浮点数精度问题
import Decimal from 'decimal.js'

// 设置精度为20位，四舍五入
Decimal.set({ precision: 20, rounding: Decimal.ROUND_HALF_UP })

// 加法
export function add(a: string | number, b: string | number): string {
  return new Decimal(a).plus(b).toString()
}

// 减法
export function sub(a: string | number, b: string | number): string {
  return new Decimal(a).minus(b).toString()
}

// 乘法
export function mul(a: string | number, b: string | number): string {
  return new Decimal(a).times(b).toString()
}

// 除法
export function div(a: string | number, b: string | number): string {
  return new Decimal(a).dividedBy(b).toString()
}

// 等于比较
export function eq(a: string | number, b: string | number): boolean {
  return new Decimal(a).equals(b)
}

// 大于比较
export function gt(a: string | number, b: string | number): boolean {
  return new Decimal(a).greaterThan(b)
}

// 大于等于比较
export function gte(a: string | number, b: string | number): boolean {
  return new Decimal(a).greaterThanOrEqualTo(b)
}

// 小于比较
export function lt(a: string | number, b: string | number): boolean {
  return new Decimal(a).lessThan(b)
}

// 小于等于比较
export function lte(a: string | number, b: string | number): boolean {
  return new Decimal(a).lessThanOrEqualTo(b)
}

// 是否为零
export function isZero(a: string | number): boolean {
  return new Decimal(a).isZero()
}

// 取绝对值
export function abs(a: string | number): string {
  return new Decimal(a).abs().toString()
}

// 四舍五入到指定小数位，默认2位
export function round(a: string | number, dp: number = 2): string {
  return new Decimal(a).toDecimalPlaces(dp).toString()
}
