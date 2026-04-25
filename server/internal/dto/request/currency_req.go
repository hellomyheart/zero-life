// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// SetCurrencyStatusReq 设置货币启用/禁用状态请求
type SetCurrencyStatusReq struct {
	IsEnabled bool `json:"is_enabled"` // 是否启用
}

// SetDefaultCurrencyReq 设置默认货币请求
type SetDefaultCurrencyReq struct {
	CurrencyID uint64 `json:"currency_id" binding:"required"` // 货币ID
}

// SetExchangeRateReq 设置汇率请求
// 设置两种货币之间的兑换汇率
type SetExchangeRateReq struct {
	FromCurrencyID uint64 `json:"from_currency_id" binding:"required"` // 源货币ID
	ToCurrencyID   uint64 `json:"to_currency_id" binding:"required"`   // 目标货币ID
	Rate           string `json:"rate" binding:"required"`             // 汇率值
}
