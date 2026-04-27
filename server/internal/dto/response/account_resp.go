// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// AccountResp 账户响应
type AccountResp struct {
	ID             uint64       `json:"id"`              // 账户ID
	Name           string       `json:"name"`            // 账户名称
	AccountNumber  string       `json:"account_number"`  // 账户号
	Type           string       `json:"type"`            // 账户类型
	CurrencyID     uint64       `json:"currency_id"`     // 货币ID
	Currency       CurrencyResp `json:"currency"`        // 关联货币信息
	InitialBalance string       `json:"initial_balance"` // 初始余额
	CurrentBalance string       `json:"current_balance"` // 当前余额（含所有交易）
	IsVirtual      bool         `json:"is_virtual"`      // 是否为虚拟账户
	Notes          string       `json:"notes"`           // 备注
	CreatedAt      time.Time    `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time    `json:"updated_at"`      // 更新时间
}
