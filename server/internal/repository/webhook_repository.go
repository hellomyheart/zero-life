package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// --- Webhook CRUD ---

func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return r.db.Create(webhook).Error
}

func (r *WebhookRepository) GetByID(id, userID uint64) (*model.Webhook, error) {
	var webhook model.Webhook
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&webhook).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *WebhookRepository) List(userID uint64) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return r.db.Save(webhook).Error
}

func (r *WebhookRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Webhook{}).Error
}

// GetByTrigger returns all active webhooks matching the given trigger for a user.
func (r *WebhookRepository) GetByTrigger(userID uint64, trigger model.WebhookTrigger) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("user_id = ? AND trigger = ? AND is_active = ?", userID, trigger, true).
		Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

// --- WebhookMessage CRUD ---

func (r *WebhookRepository) CreateMessage(msg *model.WebhookMessage) error {
	return r.db.Create(msg).Error
}

func (r *WebhookRepository) UpdateMessage(msg *model.WebhookMessage) error {
	return r.db.Save(msg).Error
}

func (r *WebhookRepository) GetMessagesByWebhookID(webhookID uint64, offset, limit int) ([]model.WebhookMessage, error) {
	var messages []model.WebhookMessage
	if err := r.db.Where("webhook_id = ?", webhookID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *WebhookRepository) CountMessagesByWebhookID(webhookID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.WebhookMessage{}).Where("webhook_id = ?", webhookID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
