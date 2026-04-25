// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateBudgetReq 创建预算请求
// 预算用于控制指定分类的支出金额
type CreateBudgetReq struct {
	Name        string   `json:"name" binding:"required"`                                // 预算名称
	Amount      string   `json:"amount" binding:"required"`                              // 预算金额，必须大于0
	Period      string   `json:"period" binding:"required,oneof=monthly yearly"`          // 预算周期：monthly(月度)/yearly(年度)
	CategoryIDs []uint64 `json:"category_ids" binding:"required,min=1"`                  // 关联分类ID列表，至少1个
}

// UpdateBudgetReq 更新预算请求
type UpdateBudgetReq struct {
	Name        string   `json:"name"`                                           // 预算名称
	Amount      string   `json:"amount"`                                         // 预算金额
	Period      string   `json:"period" binding:"omitempty,oneof=monthly yearly"` // 预算周期
	CategoryIDs []uint64 `json:"category_ids"`                                   // 关联分类ID列表
	IsEnabled   *bool    `json:"is_enabled"`                                     // 是否启用
}
