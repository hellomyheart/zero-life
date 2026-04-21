package request

type CreateBudgetReq struct {
	Name        string   `json:"name" binding:"required"`
	Amount      string   `json:"amount" binding:"required"`
	Period      string   `json:"period" binding:"required,oneof=monthly yearly"`
	CategoryIDs []uint64 `json:"category_ids" binding:"required,min=1"`
}

type UpdateBudgetReq struct {
	Name        string   `json:"name"`
	Amount      string   `json:"amount"`
	Period      string   `json:"period" binding:"omitempty,oneof=monthly yearly"`
	CategoryIDs []uint64 `json:"category_ids"`
	IsEnabled   *bool    `json:"is_enabled"`
}
