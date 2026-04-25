// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateReconciliationReq 创建对账请求
// 用于创建账户对账记录，比较账面余额和实际余额
type CreateReconciliationReq struct {
	AccountID       uint64 `json:"account_id" binding:"required"`         // 对账账户ID
	StartDate       string `json:"start_date" binding:"required"`         // 对账开始日期
	EndDate         string `json:"end_date" binding:"required"`           // 对账结束日期
	StartingBalance string `json:"starting_balance" binding:"required"`   // 期初余额
	EndingBalance   string `json:"ending_balance" binding:"required"`     // 期末余额
}

// UpdateReconciliationReq 更新对账请求
type UpdateReconciliationReq struct {
	EndingBalance string `json:"ending_balance"`                            // 期末余额
	Status        string `json:"status" binding:"omitempty,oneof=open closed"` // 对账状态：open(进行中)/closed(已关闭)
}

// ReconciliationListReq 对账列表查询请求
type ReconciliationListReq struct {
	AccountID *uint64 `form:"account_id"`             // 按账户ID过滤
	Page      int     `form:"page,default=1"`         // 页码
	PageSize  int     `form:"page_size,default=20"`   // 每页数量
}