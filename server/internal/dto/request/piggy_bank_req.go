// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreatePiggyBankReq 创建存钱罐请求
type CreatePiggyBankReq struct {
	Name         string  `json:"name" binding:"required"`         // 存钱罐名称
	TargetAmount string  `json:"target_amount" binding:"required"` // 目标金额
	AccountID    uint64  `json:"account_id" binding:"required"`   // 关联账户ID，存取金额从此账户扣除/增加
	TargetDate   *string `json:"target_date"`                     // 目标达成日期（可选）
	Notes        string  `json:"notes"`                           // 备注
}

// UpdatePiggyBankReq 更新存钱罐请求
type UpdatePiggyBankReq struct {
	Name         string  `json:"name"`          // 存钱罐名称
	TargetAmount *string `json:"target_amount"` // 目标金额
	TargetDate   *string `json:"target_date"`    // 目标达成日期
	Notes        string  `json:"notes"`          // 备注
}

// AddAmountReq 存入金额请求
// 向存钱罐中增加金额
type AddAmountReq struct {
	Amount string `json:"amount" binding:"required"` // 存入金额
	Note   string `json:"note"`                       // 存入备注
}

// RemoveAmountReq 取出金额请求
// 从存钱罐中减少金额
type RemoveAmountReq struct {
	Amount string `json:"amount" binding:"required"` // 取出金额
	Note   string `json:"note"`                      // 取出备注
}
