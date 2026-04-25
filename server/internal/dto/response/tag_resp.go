// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// TagResp 标签响应
type TagResp struct {
	ID              uint64    `json:"id"`               // 标签ID
	Name            string    `json:"name"`             // 标签名称
	Color           string    `json:"color"`            // 标签颜色
	TransactionCount int64    `json:"transaction_count"` // 关联交易数量
	CreatedAt       time.Time `json:"created_at"`       // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`       // 更新时间
}
