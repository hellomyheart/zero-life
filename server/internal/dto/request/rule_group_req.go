package request

type CreateRuleGroupReq struct {
	Name     string `json:"name" binding:"required"`
	Order    int    `json:"order"`
	IsActive *bool  `json:"is_active"`
}

type UpdateRuleGroupReq struct {
	Name     string `json:"name"`
	Order    *int   `json:"order"`
	IsActive *bool  `json:"is_active"`
}

type ExecuteRuleGroupReq struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
