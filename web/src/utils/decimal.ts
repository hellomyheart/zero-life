import Decimal from 'decimal.js'

Decimal.set({ precision: 20, rounding: Decimal.ROUND_HALF_UP })

export function add(a: string | number, b: string | number): string {
  return new Decimal(a).plus(b).toString()
}

export function sub(a: string | number, b: string | number): string {
  return new Decimal(a).minus(b).toString()
}

export function mul(a: string | number, b: string | number): string {
  return new Decimal(a).times(b).toString()
}

export function div(a: string | number, b: string | number): string {
  return new Decimal(a).dividedBy(b).toString()
}

export function eq(a: string | number, b: string | number): boolean {
  return new Decimal(a).equals(b)
}

export function gt(a: string | number, b: string | number): boolean {
  return new Decimal(a).greaterThan(b)
}

export function gte(a: string | number, b: string | number): boolean {
  return new Decimal(a).greaterThanOrEqualTo(b)
}

export function lt(a: string | number, b: string | number): boolean {
  return new Decimal(a).lessThan(b)
}

export function lte(a: string | number, b: string | number): boolean {
  return new Decimal(a).lessThanOrEqualTo(b)
}

export function isZero(a: string | number): boolean {
  return new Decimal(a).isZero()
}

export function abs(a: string | number): string {
  return new Decimal(a).abs().toString()
}

export function round(a: string | number, dp: number = 2): string {
  return new Decimal(a).toDecimalPlaces(dp).toString()
}
