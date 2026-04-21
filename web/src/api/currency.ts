import { get, put, post } from '@/utils/request'
import type { Currency, ExchangeRate, SetCurrencyStatusReq, SetDefaultCurrencyReq, SetExchangeRateReq } from '@/types/currency'

export function list() {
  return get<Currency[]>('/currencies')
}

export function updateStatus(id: string, data: SetCurrencyStatusReq) {
  return put<Currency>(`/currencies/${id}/status`, data)
}

export function setDefault(data: SetDefaultCurrencyReq) {
  return put<Currency>('/currencies/default', data)
}

export function getExchangeRates() {
  return get<ExchangeRate[]>('/currencies/exchange-rates')
}

export function setExchangeRate(data: SetExchangeRateReq) {
  return post<ExchangeRate>('/currencies/exchange-rates', data)
}
