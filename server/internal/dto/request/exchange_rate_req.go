// Package request 定义API请求参数结构
package request

import "time"

// CreateExchangeRateReq 创建汇率请求
type CreateExchangeRateReq struct {
	FromCurrencyID uint64    `json:"from_currency_id" binding:"required"` // 源货币ID
	ToCurrencyID   uint64    `json:"to_currency_id" binding:"required"`   // 目标货币ID
	Date           time.Time `json:"date" binding:"required"`             // 汇率日期
	Rate           string    `json:"rate" binding:"required"`             // 汇率值
	UserRate       string    `json:"user_rate"`                          // 用户自定义汇率（可选）
}

// UpdateExchangeRateReq 更新汇率请求
type UpdateExchangeRateReq struct {
	Rate     string `json:"rate"`      // 汇率值
	UserRate string `json:"user_rate"` // 用户自定义汇率
}

// ExchangeRateListReq 汇率列表查询请求
type ExchangeRateListReq struct {
	FromCurrencyID uint64 `form:"from_currency_id"` // 源货币ID（可选）
	ToCurrencyID   uint64 `form:"to_currency_id"`   // 目标货币ID（可选）
	StartDate      string `form:"start_date"`       // 开始日期（可选）
	EndDate        string `form:"end_date"`         // 结束日期（可选）
	Page           int    `form:"page"`             // 页码
	PageSize       int    `form:"page_size"`        // 每页记录数
}
