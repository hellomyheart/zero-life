// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// AutocompleteItemResp 自动补全项响应
// 用于前端搜索框的下拉选项
type AutocompleteItemResp struct {
	ID     uint64 `json:"id"`                // 项目ID
	Name   string `json:"name"`              // 项目名称
	Type   string `json:"type,omitempty"`    // 项目类型（账户类型等）
	Symbol string `json:"symbol,omitempty"`  // 货币符号
	Code   string `json:"code,omitempty"`    // 货币代码
	Color  string `json:"color,omitempty"`   // 标签颜色
}
