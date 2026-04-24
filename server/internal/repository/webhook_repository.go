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
	if err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return r.db.Save(webhook).Error
}

func (r *WebhookRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("webhook_id = ?", id).Delete(&model.WebhookDelivery{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Webhook{}).Error
	})
}

func (r *WebhookRepository) GetActiveByTrigger(userID uint64, trigger model.WebhookTrigger) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("user_id = ? AND trigger = ? AND is_active = ?", userID, trigger, true).
		Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *WebhookRepository) CreateDelivery(delivery *model.WebhookDelivery) error {
	return r.db.Create(delivery).Error
}

func (r *WebhookRepository) ListDeliveries(webhookID uint64, offset, limit int) ([]model.WebhookDelivery, error) {
	var deliveries []model.WebhookDelivery
	if err := r.db.Where("webhook_id = ?", webhookID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}
