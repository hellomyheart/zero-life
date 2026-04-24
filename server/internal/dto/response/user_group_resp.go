// Package response 定义API响应数据结构
package response

import "time"

// UserGroupResp 用户组响应
type UserGroupResp struct {
	ID        uint64    `json:"id"`         // 用户组ID
	Title     string    `json:"title"`      // 用户组标题
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
