package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/pkg/webhook"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// WebhookNotifier defines the interface for triggering webhooks.
// This interface decouples TransactionService from WebhookService to avoid circular dependencies.
type WebhookNotifier interface {
	TriggerWebhooks(userID uint64, trigger model.WebhookTrigger, data interface{})
}

type WebhookService struct {
	webhookRepo *repository.WebhookRepository
}

func NewWebhookService(webhookRepo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{webhookRepo: webhookRepo}
}

// --- Webhook CRUD ---

func (s *WebhookService) Create(userID uint64, req *request.CreateWebhookReq) (*response.WebhookResp, error) {
	if err := webhook.ValidateURL(req.URL); err != nil {
		return nil, errcode.ErrWebhookURLInvalid
	}

	wh := &model.Webhook{
		UserID:  userID,
		URL:     req.URL,
		Trigger: model.WebhookTrigger(req.Trigger),
		Secret:  req.Secret,
	}

	if err := s.webhookRepo.Create(wh); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(wh), nil
}

func (s *WebhookService) Get(userID, id uint64) (*response.WebhookResp, error) {
	wh, err := s.webhookRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(wh), nil
}

func (s *WebhookService) List(userID uint64) ([]response.WebhookResp, error) {
	webhooks, err := s.webhookRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.WebhookResp, 0, len(webhooks))
	for _, wh := range webhooks {
		items = append(items, *s.toResp(&wh))
	}
	return items, nil
}

func (s *WebhookService) Update(userID, id uint64, req *request.UpdateWebhookReq) (*response.WebhookResp, error) {
	wh, err := s.webhookRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if err := webhook.ValidateURL(req.URL); err != nil {
		return nil, errcode.ErrWebhookURLInvalid
	}

	wh.URL = req.URL
	wh.Trigger = model.WebhookTrigger(req.Trigger)
	wh.Secret = req.Secret
	if req.IsActive != nil {
		wh.IsActive = *req.IsActive
	}

	if err := s.webhookRepo.Update(wh); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(wh), nil
}

func (s *WebhookService) Delete(userID, id uint64) error {
	return s.webhookRepo.Delete(id, userID)
}

// --- Trigger & Messages ---

// TriggerWebhooks finds matching webhooks and sends them asynchronously.
func (s *WebhookService) TriggerWebhooks(userID uint64, trigger model.WebhookTrigger, data interface{}) {
	webhooks, err := s.webhookRepo.GetByTrigger(userID, trigger)
	if err != nil {
		log.Printf("[WebhookService] TriggerWebhooks: failed to get webhooks for user %d trigger %s: %v", userID, trigger, err)
		return
	}

	if len(webhooks) == 0 {
		return
	}

	payload := webhook.BuildPayload(string(trigger), data)
	payloadBytes, _ := json.Marshal(payload)

	for _, wh := range webhooks {
		go func(wh model.Webhook) {
			result := webhook.Send(wh.URL, wh.Secret, payload)

			now := time.Now()
			msg := &model.WebhookMessage{
				WebhookID:    wh.ID,
				RequestBody:  string(payloadBytes),
				ResponseBody: result.ResponseBody,
				Attempts:     3,
				SentAt:       &now,
				CreatedAt:    time.Now(),
			}

			if result.Success {
				msg.Status = model.WebhookMessageSuccess
				msg.ResponseCode = &result.ResponseCode
			} else {
				msg.Status = model.WebhookMessageFailed
				if result.ResponseCode > 0 {
					msg.ResponseCode = &result.ResponseCode
				}
			}

			if err := s.webhookRepo.CreateMessage(msg); err != nil {
				log.Printf("[WebhookService] TriggerWebhooks: failed to save message for webhook %d: %v", wh.ID, err)
			}
		}(wh)
	}
}

// GetMessages returns paginated messages for a webhook.
func (s *WebhookService) GetMessages(userID, webhookID uint64, req *request.WebhookMessageListReq) (*pagination.Result, error) {
	// Verify webhook belongs to user
	_, err := s.webhookRepo.GetByID(webhookID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	messages, err := s.webhookRepo.GetMessagesByWebhookID(webhookID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.webhookRepo.CountMessagesByWebhookID(webhookID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.WebhookMessageResp, 0, len(messages))
	for _, m := range messages {
		items = append(items, s.messageToResp(&m))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *WebhookService) toResp(wh *model.Webhook) *response.WebhookResp {
	return &response.WebhookResp{
		ID:        wh.ID,
		URL:       wh.URL,
		Trigger:   string(wh.Trigger),
		Secret:    wh.Secret,
		IsActive:  wh.IsActive,
		CreatedAt: wh.CreatedAt,
		UpdatedAt: wh.UpdatedAt,
	}
}

func (s *WebhookService) messageToResp(m *model.WebhookMessage) response.WebhookMessageResp {
	return response.WebhookMessageResp{
		ID:           m.ID,
		WebhookID:    m.WebhookID,
		RequestBody:  m.RequestBody,
		ResponseCode: m.ResponseCode,
		ResponseBody: m.ResponseBody,
		Attempts:     m.Attempts,
		Status:       string(m.Status),
		SentAt:       m.SentAt,
		CreatedAt:    m.CreatedAt,
	}
}
