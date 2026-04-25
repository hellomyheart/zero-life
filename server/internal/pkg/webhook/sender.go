// Package webhook 提供Webhook消息发送和验证功能
// 支持HTTP POST方式向外部URL发送Webhook通知，包含HMAC-SHA256签名验证
// 内置重试机制（最多3次，指数退避），确保消息可靠投递
package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 重试配置常量
const (
	maxRetries   = 3                      // 最大重试次数
	initialDelay = 1 * time.Second        // 初始重试延迟（1秒），后续指数退避（2秒、4秒）
)

// SendResult Webhook发送结果
// 记录单次Webhook发送的响应状态和结果
type SendResult struct {
	ResponseCode int    // HTTP响应状态码
	ResponseBody string // HTTP响应体内容
	Success      bool   // 是否发送成功（HTTP 2xx视为成功）
}

// Send 发送Webhook通知
// 将payload以JSON格式POST到指定URL，支持HMAC-SHA256签名和重试机制
// 参数：
//   - url: Webhook接收URL
//   - secret: 签名密钥，为空则不签名
//   - payload: 要发送的数据，将被序列化为JSON
// 返回：
//   - SendResult: 发送结果（含状态码、响应体、是否成功）
// 业务规则：
//   - 最多重试3次，采用指数退避策略（1秒、2秒、4秒）
//   - HTTP 2xx响应视为成功，非2xx响应触发重试
//   - secret非空时，使用HMAC-SHA256对请求体签名，签名值放入X-Webhook-Signature头
//   - HTTP请求超时时间为10秒
func Send(url string, secret string, payload interface{}) SendResult {
	body, err := json.Marshal(payload)
	if err != nil {
		return SendResult{Success: false}
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := initialDelay * (1 << (attempt - 1)) // 1s, 2s, 4s
			time.Sleep(delay)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		// Sign with HMAC-SHA256 if secret is provided
		if secret != "" {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(body)
			sig := hex.EncodeToString(mac.Sum(nil))
			req.Header.Set("X-Webhook-Signature", sig)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return SendResult{
				ResponseCode: resp.StatusCode,
				ResponseBody: string(respBody),
				Success:      true,
			}
		}

		// Non-2xx response: retry
		if attempt == maxRetries-1 {
			return SendResult{
				ResponseCode: resp.StatusCode,
				ResponseBody: string(respBody),
				Success:      false,
			}
		}
	}

	return SendResult{Success: false}
}

// BuildPayload 构建标准Webhook负载
// 将触发事件和数据封装为标准格式，包含事件类型、数据和时间戳
// 参数：
//   - trigger: 触发事件类型（如transaction.created）
//   - data: 事件相关数据
// 返回：
//   - map[string]interface{}: 标准Webhook负载
func BuildPayload(trigger string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"event":     trigger,
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// ValidateURL 验证Webhook URL是否合法
// 业务规则：URL必须以https://开头，确保传输安全
// 参数：
//   - url: 待验证的URL
// 返回：
//   - error: URL不合法时返回错误
func ValidateURL(url string) error {
	if len(url) < 8 || url[:8] != "https://" {
		return fmt.Errorf("webhook URL must start with https://")
	}
	return nil
}
