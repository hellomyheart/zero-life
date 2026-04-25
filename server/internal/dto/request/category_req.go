// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateCategoryReq 创建分类请求
// 支持设置父分类（最多2级层级）、图标和备注
type CreateCategoryReq struct {
	Name     string  `json:"name" binding:"required"` // 分类名称，同一用户下不能重复
	ParentID *uint64 `json:"parent_id"`               // 父分类ID，为空表示顶级分类
	Icon     string  `json:"icon"`                    // 分类图标
	Notes    string  `json:"notes"`                   // 备注
}

// UpdateCategoryReq 更新分类请求
type UpdateCategoryReq struct {
	Name      string  `json:"name"`       // 分类名称
	Icon      string  `json:"icon"`       // 分类图标
	Notes     string  `json:"notes"`      // 备注
	SortOrder *int    `json:"sort_order"` // 排序序号
}
