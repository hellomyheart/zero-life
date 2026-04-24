package request

type CreateWebhookReq struct {
	Name    string `json:"name" binding:"required"`
	URL     string `json:"url" binding:"required,url"`
	Trigger string `json:"trigger" binding:"required,oneof=transaction.created transaction.updated transaction.deleted bill.paid budget.created budget.updated budget.deleted"`
}

type UpdateWebhookReq struct {
	Name     string  `json:"name"`
	URL      string  `json:"url" binding:"omitempty,url"`
	Trigger  string  `json:"trigger" binding:"omitempty,oneof=transaction.created transaction.updated transaction.deleted bill.paid budget.created budget.updated budget.deleted"`
	IsActive *bool   `json:"is_active"`
}