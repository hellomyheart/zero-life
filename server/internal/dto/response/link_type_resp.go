// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// LinkTypeResp 关联类型响应
type LinkTypeResp struct {
	ID            uint64    `json:"id"`             // 关联类型ID
	Name          string    `json:"name"`           // 关联类型名称
	Outward       string    `json:"outward"`        // 正向描述
	Inward        string    `json:"inward"`         // 反向描述
	IsDirectional bool      `json:"is_directional"` // 是否有方向性
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
}
