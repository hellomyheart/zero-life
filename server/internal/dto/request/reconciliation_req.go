package request

type GetReconciliationReq struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type SubmitReconciliationReq struct {
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
	SubmittedBalance string `json:"submitted_balance" binding:"required"`
}
