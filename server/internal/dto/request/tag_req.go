// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateTagReq 创建标签请求
type CreateTagReq struct {
	Name  string `json:"name" binding:"required"` // 标签名称，同一用户下不能重复
	Color string `json:"color"`                   // 标签颜色（十六进制色值）
}

// UpdateTagReq 更新标签请求
type UpdateTagReq struct {
	Name  string `json:"name"`  // 标签名称
	Color string `json:"color"` // 标签颜色
}
