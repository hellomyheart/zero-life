export interface Currency {
  id: string
  code: string
  name: string
  symbol: string
  is_enabled: boolean
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface ExchangeRate {
  source_currency: string
  target_currency: string
  rate: string
  updated_at: string
}

export interface SetCurrencyStatusReq {
  is_enabled: boolean
}

export interface SetDefaultCurrencyReq {
  currency_id: string
}

export interface SetExchangeRateReq {
  source_currency: string
  target_currency: string
  rate: string
}
