// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// TransactionResp 交易响应
type TransactionResp struct {
	ID            uint64        `json:"id"`                         // 交易ID
	Type          string        `json:"type"`                       // 交易类型
	Date          time.Time     `json:"date"`                       // 交易日期
	Description   string        `json:"description"`                // 交易描述
	Amount        string        `json:"amount"`                     // 交易金额
	SourceID      uint64        `json:"source_id"`                  // 源账户ID
	Source        AccountResp   `json:"source"`                     // 源账户信息
	DestinationID *uint64       `json:"destination_id"`             // 目标账户ID
	Destination   *AccountResp  `json:"destination,omitempty"`      // 目标账户信息
	CategoryID    *uint64       `json:"category_id"`                // 分类ID
	Category      *CategoryResp `json:"category,omitempty"`         // 分类信息
	Notes         string        `json:"notes"`                      // 备注
	Tags          []TagResp     `json:"tags"`                       // 标签列表
	Splits        []SplitResp   `json:"splits,omitempty"`           // 拆分项列表
	BillID        *uint64       `json:"bill_id,omitempty"`          // 关联账单ID
	CreatedAt     time.Time     `json:"created_at"`                 // 创建时间
	UpdatedAt     time.Time     `json:"updated_at"`                 // 更新时间
}

// SplitResp 交易拆分项响应
type SplitResp struct {
	ID         uint64        `json:"id"`                        // 拆分项ID
	Amount     string        `json:"amount"`                    // 拆分金额
	CategoryID *uint64       `json:"category_id"`               // 分类ID
	Category   *CategoryResp `json:"category,omitempty"`        // 分类信息
	Tags       []TagResp     `json:"tags"`                      // 标签列表
	Notes      string        `json:"notes"`                     // 备注
}
