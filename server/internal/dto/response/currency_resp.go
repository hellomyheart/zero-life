// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// CurrencyResp 货币响应
type CurrencyResp struct {
	ID            uint64    `json:"id"`             // 货币ID
	Code          string    `json:"code"`           // 货币代码（如CNY、USD）
	Name          string    `json:"name"`           // 货币名称
	Symbol        string    `json:"symbol"`         // 货币符号（如¥、$）
	DecimalPlaces int       `json:"decimal_places"` // 小数位数
	IsEnabled     bool      `json:"is_enabled"`     // 是否启用
	IsDefault     bool      `json:"is_default"`     // 是否为默认货币
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
}

// ExchangeRateResp 汇率响应
type ExchangeRateResp struct {
	ID             uint64    `json:"id"`               // 汇率ID
	FromCurrencyID uint64    `json:"from_currency_id"` // 源货币ID
	ToCurrencyID   uint64    `json:"to_currency_id"`   // 目标货币ID
	Rate           string    `json:"rate"`             // 汇率值
	UpdatedAt      time.Time `json:"updated_at"`       // 更新时间
}
