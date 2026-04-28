// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateWebhookReq 创建Webhook请求
type CreateWebhookReq struct {
	Name    string `json:"name" binding:"required"`    // Webhook名称
	URL     string `json:"url" binding:"required,url"` // Webhook接收URL，必须以https://开头
	Trigger string `json:"trigger" binding:"required,oneof=transaction.created transaction.updated transaction.deleted recurring_transaction.executed budget.created budget.updated budget.deleted"` // 触发事件类型
}

// UpdateWebhookReq 更新Webhook请求
type UpdateWebhookReq struct {
	Name     string  `json:"name"`                                                                                          // Webhook名称
	URL      string  `json:"url" binding:"omitempty,url"`                                                                   // Webhook接收URL
	Trigger  string  `json:"trigger" binding:"omitempty,oneof=transaction.created transaction.updated transaction.deleted recurring_transaction.executed budget.created budget.updated budget.deleted"` // 触发事件类型
	IsActive *bool   `json:"is_active"`                                                                                     // 是否启用
}