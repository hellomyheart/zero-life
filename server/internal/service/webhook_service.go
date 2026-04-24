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

// WebhookNotifier Webhook通知接口
// 定义了触发Webhook的方法，用于解耦TransactionService和WebhookService
type WebhookNotifier interface {
	Trigger(userID uint64, trigger model.WebhookTrigger, payload interface{})
	TriggerWebhooks(userID uint64, triggerType string, payload interface{})
}

// WebhookService Webhook服务
// 负责处理Webhook的增删改查、事件触发和HTTP投递
// 依赖webhookRepo进行Webhook数据访问
type WebhookService struct {
	webhookRepo *repository.WebhookRepository // Webhook数据访问对象
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService(webhookRepo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{webhookRepo: webhookRepo}
}

// Create 创建Webhook
// 自动生成签名密钥，默认启用
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、URL、触发类型）
// 返回：
//   - *response.WebhookResp: 创建成功的Webhook信息
//   - error: 错误信息
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

// Get 获取单个Webhook详情
// 参数：
//   - userID: 用户ID
//   - id: Webhook ID
// 返回：
//   - *response.WebhookResp: Webhook信息
//   - error: 错误信息
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

// List 获取用户所有Webhook列表
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.WebhookResp: Webhook列表
//   - error: 错误信息
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

// Update 更新Webhook信息
// 支持更新名称、URL、触发类型、启用状态
// 参数：
//   - userID: 用户ID
//   - id: Webhook ID
//   - req: 更新请求参数
// 返回：
//   - *response.WebhookResp: 更新后的Webhook信息
//   - error: 错误信息
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

// Delete 删除Webhook
// 参数：
//   - userID: 用户ID
//   - id: Webhook ID
// 返回：
//   - error: 错误信息
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

// Trigger 触发指定类型的Webhook
// 获取用户所有匹配触发类型的活跃Webhook，逐个投递
// 参数：
//   - userID: 用户ID
//   - trigger: 触发类型
//   - payload: 通知负载数据
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

// deliver 投递Webhook通知
// 计算HMAC-SHA256签名，发送HTTP POST请求，记录投递结果
// 请求头包含：Content-Type、X-Webhook-Signature（签名）、X-Webhook-Trigger（触发类型）
// 超时时间为10秒
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

// recordDelivery 记录Webhook投递结果
// 保存HTTP状态码、错误信息（如有）和投递时间
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

// TriggerWebhooks 触发Webhook（通过字符串触发类型）
// 是Trigger方法的字符串参数版本，供TransactionService等调用
func (s *WebhookService) TriggerWebhooks(userID uint64, triggerType string, payload interface{}) {
	trigger := model.WebhookTrigger(triggerType)
	s.Trigger(userID, trigger, payload)
}

// ListDeliveries 获取Webhook的投递记录列表
// 参数：
//   - userID: 用户ID
//   - webhookID: Webhook ID
//   - page: 页码
//   - pageSize: 每页数量
// 返回：
//   - []response.WebhookDeliveryResp: 投递记录列表
//   - error: 错误信息
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

// toResp 将Webhook模型转换为响应对象
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

// generateWebhookSecret 生成Webhook签名密钥
// 使用SHA256对时间戳进行哈希生成随机密钥
func generateWebhookSecret() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().UnixMilli()))))
}

// computeHMAC 计算HMAC-SHA256签名
// 用于验证Webhook投递的完整性和真实性
func computeHMAC(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return fmt.Sprintf("%x", mac.Sum(nil))
}