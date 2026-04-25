// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单、Webhook 等
package model

import (
	"time"

	"gorm.io/gorm"
)

// WebhookTrigger Webhook 触发事件类型
// 定义哪些事件可以触发 Webhook 推送
type WebhookTrigger string

const (
	// WebhookTriggerTransactionCreate 交易创建时触发
	// 当新建一笔交易时，推送通知到配置的 URL
	WebhookTriggerTransactionCreate WebhookTrigger = "TRANSACTION_CREATE"

	// WebhookTriggerTransactionUpdate 交易更新时触发
	// 当修改一笔交易时，推送通知到配置的 URL
	WebhookTriggerTransactionUpdate WebhookTrigger = "TRANSACTION_UPDATE"

	// WebhookTriggerTransactionDelete 交易删除时触发
	// 当删除一笔交易时，推送通知到配置的 URL
	WebhookTriggerTransactionDelete WebhookTrigger = "TRANSACTION_DELETE"
)

// WebhookMessageStatus Webhook 消息状态
// 定义 Webhook 推送消息的生命周期状态
type WebhookMessageStatus string

const (
	// WebhookMessagePending 待发送
	// 消息已创建，等待发送
	WebhookMessagePending WebhookMessageStatus = "pending"

	// WebhookMessageSuccess 发送成功
	// 目标服务器返回 2xx 状态码
	WebhookMessageSuccess WebhookMessageStatus = "success"

	// WebhookMessageFailed 发送失败
	// 目标服务器返回非 2xx 状态码或网络错误
	WebhookMessageFailed WebhookMessageStatus = "failed"
)

// Webhook Webhook 模型，对应 webhooks 表
//
// 功能说明：
// - 配置外部回调地址，在特定事件发生时推送通知
// - 支持签名密钥，确保推送消息的安全性
// - 支持启用/禁用，便于维护
// - 记录最后投递时间，监控推送状态
//
// 使用场景：
// - 交易同步：交易创建时推送到记账软件
// - 消息通知：交易更新时发送钉钉/飞书通知
// - 数据分析：交易数据推送到 BI 系统
// - 自动化：触发外部工作流（如报销审批）
//
// 签名验证：
// - Secret 字段用于生成 HMAC 签名
// - 推送时在 HTTP Header 中携带签名
// - 接收方通过签名验证消息来源的合法性
//
// 示例：
//   Webhook：交易同步
//   URL: https://api.example.com/webhook/transactions
//   触发事件：TRANSACTION_CREATE
//   签名密钥：my-secret-key
type Webhook struct {
	// ID Webhook 唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的 Webhook 配置，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Name Webhook 名称
	// 用户自定义的名称，便于识别
	// 如："交易同步"、"钉钉通知"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// URL 回调 URL
	// 事件触发时，系统向此 URL 发送 HTTP POST 请求
	// 如：https://api.example.com/webhook/transactions
	// gorm:"size:500" 限制最大长度为 500 个字符
	URL string `gorm:"not null;size:500" json:"url"`

	// Trigger 触发事件类型
	// 取值为 WebhookTrigger 枚举：TRANSACTION_CREATE/UPDATE/DELETE
	// gorm:"size:50" 限制最大长度为 50 个字符
	Trigger WebhookTrigger `gorm:"not null;size:50" json:"trigger"`

	// Secret 签名密钥
	// 用于生成 HMAC-SHA256 签名，确保推送消息的安全性
	// 推送时在 HTTP Header 中携带 X-Webhook-Signature
	// gorm:"size:255" 限制最大长度为 255 个字符
	Secret string `gorm:"size:255" json:"secret"`

	// IsActive 是否启用
	// true: 事件触发时推送通知
	// false: 暂停推送（保留配置）
	// 默认值：true
	IsActive bool `gorm:"default:true" json:"is_active"`

	// LastDeliveredAt 最后投递时间（可选）
	// 指针类型表示可为 nil（从未投递过）
	// 非 nil: 最后一次成功投递的时间
	// 用于监控 Webhook 是否正常工作
	LastDeliveredAt *time.Time `gorm:"index" json:"last_delivered_at,omitempty"`

	// CreatedAt Webhook 创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt Webhook 最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 Webhook 模型对应的数据库表名为 webhooks
func (Webhook) TableName() string { return "webhooks" }

// WebhookMessage Webhook 消息模型，对应 webhook_messages 表
//
// 功能说明：
// - 记录每次 Webhook 推送的请求和响应
// - 支持重试机制（通过 Attempts 记录重试次数）
// - 跟踪消息状态（待发送/成功/失败）
//
// 使用场景：
// - 排查推送失败原因
// - 统计推送成功率
// - 实现重试机制
//
// 示例：
//   消息记录：
//   Webhook ID: 1
//   请求体：{"event":"TRANSACTION_CREATE","data":{...}}
//   响应码：200
//   响应体：{"status":"ok"}
//   状态：success
type WebhookMessage struct {
	// ID 消息唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// WebhookID Webhook ID
	// 关联到 webhooks 表
	// gorm:"index" 创建索引，加速按 Webhook 查询消息
	WebhookID uint64 `gorm:"not null;index" json:"webhook_id"`

	// RequestBody 请求体
	// 推送到目标 URL 的 JSON 数据
	// gorm:"type:text" 支持较长的 JSON 内容
	RequestBody string `gorm:"type:text" json:"request_body"`

	// ResponseCode HTTP 响应状态码（可选）
	// 指针类型表示可为 nil（尚未发送或发送失败无响应）
	// 2xx: 成功
	// 4xx/5xx: 失败
	ResponseCode *int `gorm:"type:integer" json:"response_code"`

	// ResponseBody 响应体
	// 目标服务器返回的响应内容
	// gorm:"type:text" 支持较长的响应内容
	ResponseBody string `gorm:"type:text" json:"response_body"`

	// Attempts 重试次数
	// 发送失败后自动重试，记录已尝试的次数
	// 默认值：0
	Attempts int `gorm:"default:0" json:"attempts"`

	// Status 消息状态
	// 取值为 WebhookMessageStatus 枚举：pending/success/failed
	// gorm:"size:20" 限制最大长度为 20 个字符
	// 默认值：pending
	Status WebhookMessageStatus `gorm:"size:20;default:pending" json:"status"`

	// SentAt 发送时间（可选）
	// 指针类型表示可为 nil（尚未发送）
	// 非 nil: 实际发送的时间
	SentAt *time.Time `gorm:"index" json:"sent_at"`

	// CreatedAt 消息创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 WebhookMessage 模型对应的数据库表名为 webhook_messages
func (WebhookMessage) TableName() string { return "webhook_messages" }

// WebhookDelivery Webhook 投递记录模型，对应 webhook_deliveries 表
//
// 功能说明：
// - 记录每次 Webhook 推送的详细信息
// - 包含推送内容、HTTP 状态码、错误信息
// - 用于排查推送问题和监控推送状态
//
// 与 WebhookMessage 的区别：
// - WebhookMessage: 侧重于消息状态跟踪（请求/响应/重试）
// - WebhookDelivery: 侧重于投递记录（推送内容/状态码/错误）
// - 两者功能重叠，后续可能合并
//
// 示例：
//   投递记录：
//   Webhook ID: 1
//   推送内容：{"event":"TRANSACTION_CREATE","data":{...}}
//   HTTP 状态码：200
//   投递时间：2026-04-25 10:30:00
type WebhookDelivery struct {
	// ID 投递记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// WebhookID Webhook ID
	// 关联到 webhooks 表
	// gorm:"index" 创建索引，加速按 Webhook 查询投递记录
	WebhookID uint64 `gorm:"not null;index" json:"webhook_id"`

	// Payload 推送内容
	// 发送到目标 URL 的 JSON 数据
	// gorm:"type:text" 支持较长的 JSON 内容
	Payload string `gorm:"type:text" json:"payload"`

	// StatusCode HTTP 状态码（可选）
	// 指针类型表示可为 nil（尚未发送或网络错误无响应）
	// 2xx: 成功
	// 4xx/5xx: 失败
	StatusCode *int `gorm:"type:integer" json:"status_code,omitempty"`

	// ErrorMessage 错误信息（可选）
	// 指针类型表示可为 nil（无错误）
	// 非 nil: 记录推送失败的错误信息
	// 如："connection refused"、"timeout"
	ErrorMessage *string `gorm:"size:500" json:"error_message,omitempty"`

	// DeliveredAt 投递时间（可选）
	// 指针类型表示可为 nil（尚未投递）
	// 非 nil: 实际投递成功的时间
	DeliveredAt *time.Time `gorm:"index" json:"delivered_at,omitempty"`

	// CreatedAt 投递记录创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 WebhookDelivery 模型对应的数据库表名为 webhook_deliveries
func (WebhookDelivery) TableName() string { return "webhook_deliveries" }
