// Package pagination 提供分页参数处理和结果封装功能
// 统一处理分页参数的规范化、偏移量计算和分页结果构建
package pagination

import "math"

// Params 分页请求参数
// 用于接收和规范化分页查询参数
type Params struct {
	Page     int `form:"page" json:"page"`         // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}

// Normalize 规范化分页参数
// 将不合法的分页参数修正为合理值
// 业务规则：
//   - 页码小于1时修正为1
//   - 每页数量小于1时修正为20
//   - 每页数量最大不超过100，防止一次查询过多数据
func (p *Params) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset 计算数据库查询偏移量
// 根据页码和每页数量计算OFFSET值
// 返回：
//   - int: 偏移量，公式为 (Page - 1) * PageSize
func (p *Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Result 分页结果结构
// 封装分页查询的结果数据，包含数据列表和分页元信息
type Result struct {
	Items      interface{} `json:"items"`       // 当前页的数据列表
	Total      int64       `json:"total"`       // 数据总条数
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页数量
	TotalPages int         `json:"total_pages"` // 总页数，向上取整
}

// NewResult 创建分页结果
// 根据数据列表、总条数和分页参数构建分页结果
// 参数：
//   - items: 当前页的数据列表
//   - total: 数据总条数
//   - params: 分页参数
// 返回：
//   - *Result: 分页结果实例
func NewResult(items interface{}, total int64, params Params) *Result {
	return &Result{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PageSize))),
	}
}
