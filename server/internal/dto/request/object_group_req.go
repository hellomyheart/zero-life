// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateObjectGroupReq 创建对象分组请求
// 将交易、循环交易、预算或存钱罐添加到分组中
type CreateObjectGroupReq struct {
	Name           string `json:"name" binding:"required"`                                        // 分组名称
	GroupableType  string `json:"groupable_type" binding:"required,oneof=transaction bill budget piggy_bank"` // 对象类型
	GroupableID    uint64 `json:"groupable_id" binding:"required"`                                // 对象ID
}

// UpdateObjectGroupReq 更新对象分组请求
type UpdateObjectGroupReq struct {
	Name string `json:"name"` // 分组名称
}