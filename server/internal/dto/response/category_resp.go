// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// CategoryResp 分类响应
// 支持树形结构，Children字段包含子分类
type CategoryResp struct {
	ID        uint64         `json:"id"`                  // 分类ID
	Name      string         `json:"name"`                // 分类名称
	ParentID  *uint64        `json:"parent_id"`           // 父分类ID
	Icon      string         `json:"icon"`                // 分类图标
	Notes     string         `json:"notes"`               // 备注
	SortOrder int            `json:"sort_order"`          // 排序序号
	Children  []*CategoryResp `json:"children,omitempty"`  // 子分类列表
	CreatedAt time.Time      `json:"created_at"`          // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`          // 更新时间
}
