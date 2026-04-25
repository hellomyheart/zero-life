// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateLinkTypeReq 创建关联类型请求
// 关联类型定义了交易之间关联关系的语义
type CreateLinkTypeReq struct {
	Name          string `json:"name" binding:"required"`          // 关联类型名称
	Outward       string `json:"outward" binding:"required"`       // 正向描述（如"冲正了"）
	Inward        string `json:"inward" binding:"required"`        // 反向描述（如"被冲正"）
	IsDirectional bool   `json:"is_directional"`                   // 是否有方向性
}

// UpdateLinkTypeReq 更新关联类型请求
type UpdateLinkTypeReq struct {
	Name          string `json:"name"`            // 关联类型名称
	Outward       string `json:"outward"`         // 正向描述
	Inward        string `json:"inward"`          // 反向描述
	IsDirectional *bool  `json:"is_directional"`  // 是否有方向性
}

// LinkTypeListReq 关联类型列表查询请求
type LinkTypeListReq struct {
	Page     int `form:"page,default=1"`        // 页码
	PageSize int `form:"page_size,default=20"`  // 每页数量
}
