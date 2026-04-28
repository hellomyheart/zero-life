// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// TagResp 标签响应
// 支持树形结构，Children字段包含子标签
// 使用指针切片避免值语义导致深层子节点丢失
type TagResp struct {
	ID               uint64        `json:"id"`                          // 标签ID
	Name             string        `json:"name"`                        // 标签名称
	Color            string        `json:"color"`                       // 标签颜色
	ParentID         *uint64       `json:"parent_id"`                   // 父标签ID，为nil时表示顶级标签
	TransactionCount int64         `json:"transaction_count"`           // 关联交易数量
	Children         []*TagResp    `json:"children,omitempty"`          // 子标签列表
	CreatedAt        time.Time     `json:"created_at"`                  // 创建时间
	UpdatedAt        time.Time     `json:"updated_at"`                  // 更新时间
}
