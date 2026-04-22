package request

type CreatePiggyBankReq struct {
	Name         string  `json:"name" binding:"required"`
	TargetAmount string  `json:"target_amount" binding:"required"`
	AccountID    uint64  `json:"account_id" binding:"required"`
	TargetDate   *string `json:"target_date"`
	Notes        string  `json:"notes"`
}

type UpdatePiggyBankReq struct {
	Name         string  `json:"name"`
	TargetAmount *string `json:"target_amount"`
	TargetDate   *string `json:"target_date"`
	Notes        string  `json:"notes"`
}

type AddAmountReq struct {
	Amount string `json:"amount" binding:"required"`
	Note   string `json:"note"`
}

type RemoveAmountReq struct {
	Amount string `json:"amount" binding:"required"`
	Note   string `json:"note"`
}
