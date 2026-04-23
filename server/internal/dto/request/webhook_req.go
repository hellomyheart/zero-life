package request

type CreateWebhookReq struct {
	URL     string `json:"url" binding:"required"`
	Trigger string `json:"trigger" binding:"required,oneof=TRANSACTION_CREATE TRANSACTION_UPDATE TRANSACTION_DELETE"`
	Secret  string `json:"secret"`
}

type UpdateWebhookReq struct {
	URL      string `json:"url" binding:"required"`
	Trigger  string `json:"trigger" binding:"required,oneof=TRANSACTION_CREATE TRANSACTION_UPDATE TRANSACTION_DELETE"`
	Secret   string `json:"secret"`
	IsActive *bool  `json:"is_active"`
}

type WebhookMessageListReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}
