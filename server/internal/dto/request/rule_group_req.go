// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateRuleGroupReq 创建规则组请求
type CreateRuleGroupReq struct {
	Name     string `json:"name" binding:"required"` // 规则组名称
	Order    int    `json:"order"`                   // 排序序号
	IsActive *bool  `json:"is_active"`               // 是否启用
}

// UpdateRuleGroupReq 更新规则组请求
type UpdateRuleGroupReq struct {
	Name     string `json:"name"`     // 规则组名称
	Order    *int   `json:"order"`    // 排序序号
	IsActive *bool  `json:"is_active"` // 是否启用
}

// ExecuteRuleGroupReq 执行规则组请求
type ExecuteRuleGroupReq struct {
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}
