package response

import "time"

type CurrencyResp struct {
	ID            uint64    `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Symbol        string    `json:"symbol"`
	DecimalPlaces int       `json:"decimal_places"`
	IsEnabled     bool      `json:"is_enabled"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ExchangeRateResp struct {
	ID             uint64    `json:"id"`
	FromCurrencyID uint64    `json:"from_currency_id"`
	ToCurrencyID   uint64    `json:"to_currency_id"`
	Rate           string    `json:"rate"`
	UpdatedAt      time.Time `json:"updated_at"`
}
