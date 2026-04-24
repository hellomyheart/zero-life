// Package response 定义API响应数据结构
package response

import "time"

// ExchangeRateResp 汇率响应
type ExchangeRateResp struct {
	ID             uint64       `json:"id"`              // 汇率ID
	UserID         uint64       `json:"user_id"`         // 用户ID
	FromCurrencyID uint64       `json:"from_currency_id"` // 源货币ID
	ToCurrencyID   uint64       `json:"to_currency_id"`   // 目标货币ID
	Date           string       `json:"date"`             // 汇率日期
	Rate           string       `json:"rate"`             // 汇率值
	UserRate       string       `json:"user_rate,omitempty"` // 用户自定义汇率
	FromCurrency   CurrencyResp `json:"from_currency,omitempty"` // 源货币信息
	ToCurrency     CurrencyResp `json:"to_currency,omitempty"`   // 目标货币信息
	CreatedAt      time.Time    `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time    `json:"updated_at"`      // 更新时间
}
