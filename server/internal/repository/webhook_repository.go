package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// WebhookRepository Webhook 仓库，负责 Webhook 配置和投递记录的数据访问。
// Webhook 允许用户在特定事件（如交易创建、账单到期）触发时，自动向指定 URL 发送 HTTP 请求。
// 每个 Webhook 包含投递记录（WebhookDelivery），记录每次发送的结果和状态。
type WebhookRepository struct {
	db *gorm.DB
}

// NewWebhookRepository 创建 Webhook 仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// Create 创建一条新的 Webhook 配置记录。
// 执行 SQL: INSERT INTO webhooks (...)
// 参数 webhook: 要创建的 Webhook 对象。
// 返回: 创建失败时返回错误。
func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return r.db.Create(webhook).Error
}

// GetByID 根据 ID 和用户 ID 获取单条 Webhook 配置。
// 执行 SQL: SELECT * FROM webhooks WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: Webhook ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的 Webhook 对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *WebhookRepository) GetByID(id, userID uint64) (*model.Webhook, error) {
	var webhook model.Webhook
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&webhook).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

// List 获取指定用户的所有 Webhook 配置，按 ID 升序排列。
// 执行 SQL: SELECT * FROM webhooks WHERE user_id = ? ORDER BY id ASC
// 参数 userID: 用户 ID。
// 返回: Webhook 列表。
func (r *WebhookRepository) List(userID uint64) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

// Update 更新 Webhook 配置。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE webhooks SET ... WHERE id = ?
// 参数 webhook: 要更新的 Webhook 对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return r.db.Save(webhook).Error
}

// Delete 删除 Webhook 及其关联的所有投递记录。
// 使用数据库事务确保原子性：先删除投递记录，再删除 Webhook 本身。
// 执行 SQL（事务内）:
//   1. DELETE FROM webhook_deliveries WHERE webhook_id = ?
//   2. DELETE FROM webhooks WHERE id = ? AND user_id = ?
// 参数 id: Webhook ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的 Webhook。
// 返回: 删除失败时返回错误。
func (r *WebhookRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("webhook_id = ?", id).Delete(&model.WebhookDelivery{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Webhook{}).Error
	})
}

// GetActiveByTrigger 获取指定用户中，匹配指定触发器且处于激活状态的 Webhook 列表。
// 当某个事件发生时（如交易创建），系统调用此方法找到需要通知的 Webhook。
// 执行 SQL: SELECT * FROM webhooks WHERE user_id = ? AND trigger = ? AND is_active = true
// 参数 userID: 用户 ID。
// 参数 trigger: Webhook 触发器类型（如 "transaction.created"）。
// 返回: 匹配的激活 Webhook 列表。
func (r *WebhookRepository) GetActiveByTrigger(userID uint64, trigger model.WebhookTrigger) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("user_id = ? AND trigger = ? AND is_active = ?", userID, trigger, true).
		Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

// CreateDelivery 创建一条 Webhook 投递记录，记录每次发送 HTTP 请求的结果。
// 执行 SQL: INSERT INTO webhook_deliveries (...)
// 参数 delivery: 要创建的投递记录对象。
// 返回: 创建失败时返回错误。
func (r *WebhookRepository) CreateDelivery(delivery *model.WebhookDelivery) error {
	return r.db.Create(delivery).Error
}

// ListDeliveries 分页获取指定 Webhook 的投递记录，按创建时间倒序排列（最新的在前）。
// 执行 SQL: SELECT * FROM webhook_deliveries WHERE webhook_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
// 参数 webhookID: Webhook ID。
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 投递记录列表。
func (r *WebhookRepository) ListDeliveries(webhookID uint64, offset, limit int) ([]model.WebhookDelivery, error) {
	var deliveries []model.WebhookDelivery
	if err := r.db.Where("webhook_id = ?", webhookID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}
