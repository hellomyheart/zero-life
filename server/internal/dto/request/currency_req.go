package request

type SetCurrencyStatusReq struct {
	IsEnabled bool `json:"is_enabled"`
}

type SetDefaultCurrencyReq struct {
	CurrencyID uint64 `json:"currency_id" binding:"required"`
}

type SetExchangeRateReq struct {
	FromCurrencyID uint64 `json:"from_currency_id" binding:"required"`
	ToCurrencyID   uint64 `json:"to_currency_id" binding:"required"`
	Rate           string `json:"rate" binding:"required"`
}
