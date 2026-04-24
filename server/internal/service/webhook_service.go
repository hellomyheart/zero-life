// Package service 业务逻辑层，实现核心业务逻辑
// WebhookService Webhook业务逻辑，处理事件通知、签名和HTTP发送
package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type WebhookNotifier interface {
	Trigger(userID uint64, trigger model.WebhookTrigger, payload interface{})
	TriggerWebhooks(userID uint64, triggerType string, payload interface{})
}

type WebhookService struct {
	webhookRepo *repository.WebhookRepository
}

func NewWebhookService(webhookRepo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{webhookRepo: webhookRepo}
}

func (s *WebhookService) Create(userID uint64, req *request.CreateWebhookReq) (*response.WebhookResp, error) {
	webhook := &model.Webhook{
		UserID:   userID,
		Name:     req.Name,
		URL:      req.URL,
		Trigger:  model.WebhookTrigger(req.Trigger),
		IsActive: true,
		Secret:   generateWebhookSecret(),
	}

	if err := s.webhookRepo.Create(webhook); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(webhook), nil
}

func (s *WebhookService) Get(userID, id uint64) (*response.WebhookResp, error) {
	webhook, err := s.webhookRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(webhook), nil
}

func (s *WebhookService) List(userID uint64) ([]response.WebhookResp, error) {
	webhooks, err := s.webhookRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.WebhookResp, 0, len(webhooks))
	for _, w := range webhooks {
		items = append(items, *s.toResp(&w))
	}
	return items, nil
}

func (s *WebhookService) Update(userID, id uint64, req *request.UpdateWebhookReq) (*response.WebhookResp, error) {
	webhook, err := s.webhookRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		webhook.Name = req.Name
	}
	if req.URL != "" {
		webhook.URL = req.URL
	}
	if req.Trigger != "" {
		webhook.Trigger = model.WebhookTrigger(req.Trigger)
	}
	if req.IsActive != nil {
		webhook.IsActive = *req.IsActive
	}

	if err := s.webhookRepo.Update(webhook); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(webhook), nil
}

func (s *WebhookService) Delete(userID, id uint64) error {
	_, err := s.webhookRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.webhookRepo.Delete(id, userID)
}

func (s *WebhookService) Trigger(userID uint64, trigger model.WebhookTrigger, payload interface{}) {
	webhooks, err := s.webhookRepo.GetActiveByTrigger(userID, trigger)
	if err != nil {
		return
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, webhook := range webhooks {
		s.deliver(&webhook, payloadBytes)
	}
}

func (s *WebhookService) deliver(webhook *model.Webhook, payload []byte) {
	signature := computeHMAC(webhook.Secret, payload)

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewReader(payload))
	if err != nil {
		s.recordDelivery(webhook.ID, payload, nil, fmt.Sprintf("failed to create request: %v", err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Trigger", string(webhook.Trigger))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.recordDelivery(webhook.ID, payload, nil, fmt.Sprintf("failed to send: %v", err))
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	var errorMsg *string
	if statusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		msg := fmt.Sprintf("HTTP %d: %s", statusCode, string(body))
		errorMsg = &msg
	}

	now := time.Now()
	webhook.LastDeliveredAt = &now
	s.webhookRepo.Update(webhook)

	s.recordDelivery(webhook.ID, payload, &statusCode, func() string {
		if errorMsg != nil {
			return *errorMsg
		}
		return ""
	}())
}

func (s *WebhookService) recordDelivery(webhookID uint64, payload []byte, statusCode *int, errorMsg string) {
	delivery := &model.WebhookDelivery{
		WebhookID:   webhookID,
		Payload:     string(payload),
		StatusCode:  statusCode,
		DeliveredAt: func() *time.Time { t := time.Now(); return &t }(),
	}
	if errorMsg != "" {
		delivery.ErrorMessage = &errorMsg
	}
	s.webhookRepo.CreateDelivery(delivery)
}

func (s *WebhookService) TriggerWebhooks(userID uint64, triggerType string, payload interface{}) {
	trigger := model.WebhookTrigger(triggerType)
	s.Trigger(userID, trigger, payload)
}

func (s *WebhookService) ListDeliveries(userID, webhookID uint64, page, pageSize int) ([]response.WebhookDeliveryResp, error) {
	_, err := s.webhookRepo.GetByID(webhookID, userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	offset := (page - 1) * pageSize
	deliveries, err := s.webhookRepo.ListDeliveries(webhookID, offset, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.WebhookDeliveryResp, 0, len(deliveries))
	for _, d := range deliveries {
		items = append(items, response.WebhookDeliveryResp{
			ID:           d.ID,
			WebhookID:    d.WebhookID,
			StatusCode:   d.StatusCode,
			ErrorMessage: d.ErrorMessage,
			DeliveredAt:  d.DeliveredAt,
			CreatedAt:    d.CreatedAt,
		})
	}
	return items, nil
}

func (s *WebhookService) toResp(w *model.Webhook) *response.WebhookResp {
	return &response.WebhookResp{
		ID:              w.ID,
		Name:            w.Name,
		URL:             w.URL,
		Trigger:         string(w.Trigger),
		IsActive:        w.IsActive,
		LastDeliveredAt: w.LastDeliveredAt,
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}

func generateWebhookSecret() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().UnixMilli()))))
}

func computeHMAC(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return fmt.Sprintf("%x", mac.Sum(nil))
}