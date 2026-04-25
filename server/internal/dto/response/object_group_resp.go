// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// ObjectGroupResp 对象分组响应
type ObjectGroupResp struct {
	ID            uint64    `json:"id"`              // 分组ID
	Name          string    `json:"name"`            // 分组名称
	GroupableType string    `json:"groupable_type"`  // 对象类型
	GroupableID   uint64    `json:"groupable_id"`    // 对象ID
	CreatedAt     time.Time `json:"created_at"`      // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`      // 更新时间
}