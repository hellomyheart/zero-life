// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateBillReq 创建账单请求
// 账单用于记录周期性支出（如房租、水电费等）
type CreateBillReq struct {
	Name       string  `json:"name" binding:"required"`                                      // 账单名称
	Amount     string  `json:"amount" binding:"required"`                                    // 账单金额，必须大于0
	RepeatRule string  `json:"repeat_rule" binding:"required,oneof=daily weekly monthly yearly"` // 重复规则：daily/weekly/monthly/yearly
	NextDue    string  `json:"next_due" binding:"required"`                                  // 下次到期日期
	SourceID   *uint64 `json:"source_id"`                                                    // 支出账户ID
	CategoryID *uint64 `json:"category_id"`                                                  // 分类ID
	Notes      string  `json:"notes"`                                                        // 备注
}

// UpdateBillReq 更新账单请求
type UpdateBillReq struct {
	Name       string  `json:"name"`                                                    // 账单名称
	Amount     string  `json:"amount"`                                                  // 账单金额
	RepeatRule string  `json:"repeat_rule" binding:"omitempty,oneof=daily weekly monthly yearly"` // 重复规则
	NextDue    string  `json:"next_due"`                                                // 下次到期日期
	SourceID   *uint64 `json:"source_id"`                                               // 支出账户ID
	CategoryID *uint64 `json:"category_id"`                                             // 分类ID
	Notes      string  `json:"notes"`                                                   // 备注
}
