package request

type CreateAccountReq struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=asset expense revenue liability"`
	CurrencyID     uint64 `json:"currency_id" binding:"required"`
	InitialBalance string `json:"initial_balance" binding:"required"`
	Notes          string `json:"notes"`
	IsVirtual      bool   `json:"is_virtual"`
}

type UpdateAccountReq struct {
	Name      string `json:"name"`
	Notes     string `json:"notes"`
	IsVirtual *bool  `json:"is_virtual"`
}

type AccountListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Type     string `form:"type"`
	Search   string `form:"search"`
	Sort     string `form:"sort,default=name"`
}
