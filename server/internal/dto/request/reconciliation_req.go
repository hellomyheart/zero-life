package request

type CreateReconciliationReq struct {
	AccountID       uint64 `json:"account_id" binding:"required"`
	StartDate       string `json:"start_date" binding:"required"`
	EndDate         string `json:"end_date" binding:"required"`
	StartingBalance string `json:"starting_balance" binding:"required"`
	EndingBalance   string `json:"ending_balance" binding:"required"`
}

type UpdateReconciliationReq struct {
	EndingBalance string `json:"ending_balance"`
	Status        string `json:"status" binding:"omitempty,oneof=open closed"`
}

type ReconciliationListReq struct {
	AccountID *uint64 `form:"account_id"`
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=20"`
}